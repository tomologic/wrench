#!/usr/bin/env bats

setup () {
    BATS_TMP_DIR=$(mktemp -d .wrench-bats.XXXXX)

    cp -r "$BATS_TEST_DIRNAME/../examples/simple" "$BATS_TMP_DIR/origin"

    pushd $BATS_TMP_DIR
    pushd origin

    # Create a bogus local origin for simple
    git init
    git config user.name "Your Name"
    git config user.email "you@example.com"
    git add .
    git commit -m "Initial commit"

    popd

    # Local clone from our local origin
    git clone ./origin ./simple
    cd simple
    git config user.name "Your Name"
    git config user.email "you@example.com"
}

teardown () {
    popd
    rm -rf $BATS_TMP_DIR

    # Remove cached test images
    docker rmi $(docker images | grep example/simple | tr -s ' ' | cut -d' ' -f3) || true
}

@test "BUMP: initial release" {
    # make sure snapshot image exists
    wrench build

    # test to bump initial version
    run wrench bump minor
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
    [ "$output" = "Released v0.1.0" ]

    # verify tag exists origin
    run cd ../origin && git rev-list v0.1.0..
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]

    # verify image exists
    run docker inspect example/simple:v0.1.0
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
}

@test "BUMP: missing docker snapshot image" {
    run wrench bump minor
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 1 ]
    [[ "$output" =~ Docker\ image\ example/simple:v0\.0\.0\-1\-.*\ does\ not\ exists ]]
}

@test "BUMP: current revision already released" {
    git tag -a v0.1.0 -m "Release v0.1.0"

    run wrench bump minor
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
    [ "$output" = "Revision already release 'v0.1.0'. Doing nothing." ]
}

@test "BUMP: release new major" {
    git tag -a v1.2.3 -m "Release v1.2.3"
    git commit -m "New feature" --allow-empty
    # make sure snapshot image exists
    wrench build

    # test to bump major version
    run wrench bump major
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
    [ "$output" = "Released v2.0.0" ]

    # verify tag exists origin
    run cd ../origin && git rev-list v2.0.0..
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]

    # verify image exists
    run docker inspect example/simple:v2.0.0
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
}

@test "BUMP: release new minor" {
    git tag -a v1.2.3 -m "Release v1.2.3"
    git commit -m "New feature" --allow-empty
    # make sure snapshot image exists
    wrench build

    # test to bump minor version
    run wrench bump minor
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
    [ "$output" = "Released v1.3.0" ]

    # verify tag exists origin
    run cd ../origin && git rev-list v1.3.0..
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]

    # verify image exists
    run docker inspect example/simple:v1.3.0
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
}

@test "BUMP: release new patch" {
    git tag -a v1.2.33 -m "Release v1.2.33"
    git commit -m "New feature" --allow-empty
    # make sure snapshot image exists
    wrench build

    # test to bump patch version
    run wrench bump patch
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
    [ "$output" = "Released v1.2.34" ]

    # verify tag exists origin
    run cd ../origin && git rev-list v1.2.34..
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]

    # verify image exists
    run docker inspect example/simple:v1.2.34
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
}

@test "BUMP: fail pushing git tag" {
    # make sure snapshot image exists
    wrench build

    # remove origin so push of tag fails
    rm -rf ../origin

    # test to bump patch version
    run wrench bump minor
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 1 ]

    # verify tag does not exists locally
    run git rev-list v1.2.34..
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 128 ]

    # verify image does not exists locally
    run docker inspect example/simple:v1.2.34
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 1 ]
}
