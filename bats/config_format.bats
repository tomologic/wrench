#!/usr/bin/env bats

setup () {
    BATS_TMP_DIR=$(mktemp -d .wrench-bats.XXXXX)
    pushd $BATS_TMP_DIR

    git init
    git config user.name "Your Name"
    git config user.email "you@example.com"

    touch a
    git add a
    git commit -m "add a"
    git tag -a v0.10.0 -m "release v0.10.0"

    cat > wrench.yml << EOF
Project:
  Name: hello
  Organization: world
Run:
  syntax: echo foo
  integration:
    Cmd: run test
    Env:
     - some=var
EOF
}

teardown () {
    popd
    rm -rf $BATS_TMP_DIR
}

@test "CONFIG: format project name" {
    ret=0
    out=$(wrench config --format '{{.Project.Name}}') || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo "out=$out"
    [ "$out" == "hello" ]
}

@test "CONFIG: format project image" {
    ret=0
    out=$(wrench config --format '{{.Project.Image}}') || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo "out=$out"
    [ "$out" == "world/hello:v0.10.0" ]
}
