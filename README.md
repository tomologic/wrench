# Wrench

[![Travis](https://img.shields.io/travis/tomologic/wrench.svg?style=flat-square)](https://travis-ci.org/tomologic/wrench)

## What is this?

This is our current try to create standards for building, testing and deploying our applications.

What we want:
- General pipelines that can build, test and deploy all our apps/services
- Useful development tooling for local development

Putting the logic in the app repos would create a lot of duplications, maintenance hell and big overhead when creating new services. Putting the logic in the pipelines will make local development difficult since all logic for building/testing/deploying exists in the CI systems. Logic for building/testing/deploying apps/services needs to be accessible for both developers and CI system.

Golang was chosen for Wrench since we can have automatic cross-compiled binaries for Linux and OSX.

Wrench is provided as [docker image](https://registry.hub.docker.com/u/tomologic/wrench/) and through our [homebrew tap](https://github.com/tomologic/homebrew-tap) for OSX users.

```
brew tap tomologic/homebrew-tap
brew install wrench
```

## Project config

### Print wrench config

To view the config that wrench will use for a project simply run the config subcommand.

```
$ cd examples/simple
$ wrench config
Project:
  Organization: example
  Name: simple
  Version: v1.0.0
Run:
  syntax-test: |
    #!/bin/bash -xe

    echo "run syntax check"
```

To get specific values use the format flag. Format will be executed as a [golang template](http://golang.org/pkg/text/template/).

```
$ cd examples/simple
$ wrench config --format '{{.Project.Name}}'
simple
```

### Wrench.yml file

Wrench will try to detect project config automatically.

- Organization is derived from current hostname.
- Name is derived from current directory.
- Version is derived from latest semver git tag.

It's possible to override all config with a _wrench.yml_ file.

```
$ cat wrench.yml
Project:
  Organization: example
  Name: real-app-name
  Version: v1.0.0
```

The _wrench.yml_ is treated by wrench as a [golang template](http://golang.org/pkg/text/template/) file where the environment is accessible through _.Environ_.

```
$ cat wrench.yml
Project:
  Organization: {{ or (.Environ.DOCKER_REGISTRY) "localhost" }}/example
  Name: {{ or (.Environ.IMAGE_NAME_PREFIX) "" }}real-app-name
  Version: v1.0.0
```

## Build

There are 2 ways that wrench can build docker images. Simple mode which is the normal docker approach and the builder which is uses a separate image to build the application artifact _(very useful for golang)_.

The mode is chosen depending on the existing Dockerfiles in current directory. Builder mode will override simple though, which makes it possible to use [Automated Builds on Docker Hub](https://docs.docker.com/docker-hub/builds/) with the regular Dockerfile and the builder image for usage with wrench.

Wrench will by default never rebuild an image if it already exists. Use rebuild flag to force a rebuild.

```
$ wrench build --rebuild
```

### VERSION environment

On build and bump wrench will add a VERSION environment to docker images. This could be utilized by the application to report it's version through an api.

It should be noted that the extra slice created during bump will make the hash of the snapshot and release docker images hash be different.

```
$ cd examples/simple
$ wrench build
INFO: Found Dockerfile, building image example/simple:v1.0.0
...
Successfully built 5eb975a7956b
$ wrench run echo-version
1.0.0
```

### Simple

Simple mode will build an image named and tagged based on project config _(either automatically detected or provided through wrench.yml)_.

```
$ cd examples/simple
$ wrench build
INFO: Found Dockerfile, building image example/simple:v1.0.0
...
Successfully built 5eb975a7956b
```

### Builder

Builder mode will use a builder image to build the final image. Builder mode will be used if a _Dockerfile.builder_ file exists.

1. Wrench builds the builder image
2. Wrench assumes that when the builder is run it will output a docker image context to stdout
3. Wrench builds the final image from the docker image context produced by the builder image

Golang example is provided in this repo:

```
$ cat examples/builder/Dockerfile.builder
```

### Test

Wrench will build a test image incase a _Dockerfile.test_ file exists. Wrench expects the first row to starts with _FROM_ and will replace the row to make it base of the application image.

Following example:

```
$ cat examples/test/Dockerfile.test
FROM
RUN pip install -r requirements-test.txt
```

The final application image might be unpractical to run tests in incase builder mode is used. You can simply tell wrench to base the test image of the builder image instead by making sure the FROM line ends with _"builder"_.

```
$ cat examples/builder/Dockerfile.test
FROM builder
WORKDIR /src
```

## Run commands

Wrench provides a subcommand to run commands inside the produced docker images that are provided in the wrench file.

Wrench will run the command in the test image if the project has one, otherwise in the final application image.

```
$ cd examples/test/
$ cat wrench.yml
Project:
  Organization: example
  Name: test
  Version: v1.0.0
Run:
  syntax-test: flake8 -v .
$ wrench run syntax-test
INFO: running syntax-test in image example/test:v1.0.0-test
directory .
checking ./server.py
```

### Environment variables

Environmental variables can be provided for run commands through the _wrench.yml_ file.

```
Project:
  Organization: example
  Name: simple
Run:
  onerow: echo "onerow"
  passthrough:
    Cmd: echo $FOO
    Env:
      - FOO={{ or (.Environ.FOO) "default-value" }}
  expanded:
    Cmd: echo "expanded"
    Env:
      - FOO=BAR
      - HELLO=WORLD
```

## Bump

Subcommand for bumping version of project. This is higly opiniated and will not work if following assumptions are not meet.

- Fast-forward only _(no merge commits used in project)_
- Snapshot image exist on host _(either already built and tested on host or pull)_
- Git access to origin _(able to push git tag)_

Example commands:

```
wrench bump major
wrench bump minor
wrench bump patch
```

Wrench will do following:

1. Generate new release version _(major, minor or patch)_
2. Git tag local git tree with new release version
3. Retag docker snapshot image to release version
4. Update VERSION env variable in release image
5. Push git tag to origin

On failure wrench will try to backtrack gracefully _(remove retagged docker image, local git tag)_.

## Push

Wrench provides a subcommand to simplify pushing of projects docker images to docker registries.

Following blocks are equivalent:

```
# Prepare tags for pushing (force required for existing tags)
docker tag example/foobar:$(git describe) registry.local:5000/example/foobar:$(git describe)
docker tag -f example/foobar:$(git describe) registry.local:5000/example/foobar:latest
docker tag -f example/foobar:$(git describe) registry.local:5000/example/foobar:prod

# Push both tags
docker push registry.local:5000/example/foobar:$(git describe)
docker push registry.local:5000/example/foobar:latest
docker push registry.local:5000/example/foobar:prod

# Remove temporary tags
docker rmi registry.local:5000/example/foobar:$(git describe)
docker rmi registry.local:5000/example/foobar:latest
docker rmi registry.local:5000/example/foobar:prod
```

Wrench will know which images to push based on project state _(release, snapshot image)_.

```
wrench push registry.local:5000 --additional-tags latest,prod
```
