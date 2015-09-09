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

    cat > wrench.yml << EOF
Project:
  Name: hello
  Organization: world
  Version: v13.1.0
  Image: custom-name
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

@test "CONFIG: configured project name" {
    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo "out=$out"
    echo $out | grep "Name: hello"
}

@test "CONFIG: configured project organization" {
    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo "out=$out"
    echo $out | grep "Organization: world"
}

@test "CONFIG: configured project version" {
    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo "out=$out"
    echo $out | grep "Version: v13.1.0"
}

@test "CONFIG: configured project image" {
    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo "out=$out"
    echo $out | grep "Image: custom-name"
}

@test "CONFIG: configured simple run" {
    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo "out=$out"
    echo $out | grep "syntax:.*"
    echo $out | grep "Cmd: echo foo"
}

@test "CONFIG: configured expanded run" {
    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo "out=$out"
    echo $out | grep "integration:.*"
    echo $out | grep "Cmd: run test.*some=var"
}
