#!/usr/bin/env bats

setup () {
    BATS_TMP_DIR=$(mktemp -d .wrench-bats.XXXXX)

    cp -r "$BATS_TEST_DIRNAME/../examples/test" "$BATS_TMP_DIR/origin"

    pushd $BATS_TMP_DIR
    pushd origin

    # Create a bogus local origin for test
    git init
    git config user.name "Your Name"
    git config user.email "you@example.com"
    git add .
    git commit -m "Initial commit"
    git tag -a v0.1.0 -m "Initial release"

    popd

    # Local clone from our local origin
    git clone ./origin ./test
    cd test
    git config user.name "Your Name"
    git config user.email "you@example.com"

    # Remove cached test images
    docker rmi example/test:v0.1.0 || true
}

teardown () {
    popd
    rm -rf $BATS_TMP_DIR

    # Remove cached test images
    docker rmi example/test:v0.1.0 || true
}

@test "EXAMPLE: build test" {
    ret=0
    out=$(wrench build) || ret=$?

    echo "out=$out"
    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo $out | grep "building.*example\/test:v0\.1\.0"
}

@test "EXAMPLE: run syntax-tests test" {
    # Build image so it already exists
    wrench build

    ret=0
    out=$(wrench run syntax-test) || ret=$?

    echo "out=$out"
    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    echo $out | grep "running syntax-test in image example/test:v0.1.0-test"
}
