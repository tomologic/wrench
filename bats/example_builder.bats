#!/usr/bin/env bats

setup () {
    BATS_TMP_DIR=$(mktemp -d .wrench-bats.XXXXX)

    cp -r "$BATS_TEST_DIRNAME/../examples/builder" "$BATS_TMP_DIR/origin"

    pushd $BATS_TMP_DIR
    pushd origin

    # Create a bogus local origin for builder
    git init
    git config user.name "Your Name"
    git config user.email "you@example.com"
    git add .
    git commit -m "Initial commit"
    git tag -a v0.1.0 -m "Initial release"

    popd

    # Local clone from our local origin
    git clone ./origin ./builder
    cd builder
    git config user.name "Your Name"
    git config user.email "you@example.com"

    # Remove cached builder images
    docker rmi example/builder:v0.1.0 || true
}

teardown () {
    popd
    rm -rf $BATS_TMP_DIR

    # Remove cached builder images
    docker rmi example/builder:v0.1.0 || true
    docker rmi example/builder:v0.1.0-test || true
    docker rmi example/builder:v0.1.0-builder || true
}

@test "EXAMPLE: build builder" {
    # Skip builder examples when docker version is not 1.7*
    # Bug in docker cli version 1.8 and 1.9 that restrict users from
    # sending docker context through stdin
    # https://github.com/docker/docker/issues/15785
    docker version | grep -q "Client version: 1.7" || skip

    run wrench build

    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]

    echo $output | grep "building.*example\/builder:v0\.1\.0"
}

@test "EXAMPLE: run syntax-tests builder" {
    # Skip builder examples when docker version is not 1.7*
    # Bug in docker cli version 1.8 and 1.9 that restrict users from
    # sending docker context through stdin
    # https://github.com/docker/docker/issues/15785
    docker version | grep -q "Client version: 1.7" || skip

    # Build image so it already exists
    wrench build

    run wrench run go-test

    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]

    echo $output | grep "PASS"
}
