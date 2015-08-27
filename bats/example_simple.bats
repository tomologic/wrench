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
    git tag -a v0.1.0 -m "Initial release"

    popd

    # Local clone from our local origin
    git clone ./origin ./simple
    cd simple
    git config user.name "Your Name"
    git config user.email "you@example.com"

    # Remove cached test images
    docker rmi example/simple:v0.1.0 || true
}

teardown () {
    popd
    rm -rf $BATS_TMP_DIR

    # Remove cached test images
    docker rmi example/simple:v0.1.0 || true
}

@test "EXAMPLE: build simple" {
    ret=0
    out=$(wrench build) || ret=$?

    echo "out=$out"
    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo $out | grep "building.*example\/simple:v0\.1\.0"
}

@test "EXAMPLE: run syntax-tests simple" {
    # Build image so it already exists
    wrench build

    ret=0
    out=$(wrench run syntax-test) || ret=$?

    echo "out=$out"
    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo $out | grep "run syntax check"
}
