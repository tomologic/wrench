#!/usr/bin/env bats

setup () {
    BATS_TMP_DIR=$(mktemp -d .wrench-bats.XXXXX)

    cp -r "$BATS_TEST_DIRNAME/../examples/simple" "$BATS_TMP_DIR/origin"

    pushd $BATS_TMP_DIR
    pushd origin

    # Add Run target
    echo '  current_version: echo "CURRENT VERSION $VERSION"' >> wrench.yml

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

    # Remove test images
    test_images=$(docker images | grep -o "example\/simple[[:space:]]*\S*" | tr -s ' ' | sed 's/\ /:/')
    if [ -n "$test_images" ]; then
        docker rmi $test_images
    fi
}

teardown () {
    popd
    rm -rf $BATS_TMP_DIR

    # Remove test images
    test_images=$(docker images | grep -o "example\/simple[[:space:]]*\S*" | tr -s ' ' | sed 's/\ /:/')
    if [ -n "$test_images" ]; then
        docker rmi $test_images
    fi
}

@test "VERSION: add version on snapshot image" {
    git commit -m "snapshot version" --allow-empty

    run wrench build
    echo "status=$status"
    echo "output=$output"
    [ "$status" -eq 0 ]

    run wrench run current_version
    echo "status=$status"
    echo "lines[1]=${lines[1]}"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" == *"CURRENT VERSION $(git describe)"* ]]
}

@test "VERSION: add version on release image" {
    run wrench build
    echo "status=$status"
    echo "output=$output"
    [ "$status" -eq 0 ]

    run wrench run current_version
    echo "status=$status"
    echo "lines[1]=${lines[1]}"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" == *"CURRENT VERSION v0.1.0"* ]]
}

@test "VERSION: update version on bump" {
    # Build snapshot version
    git commit -m "snapshot version" --allow-empty
    run wrench build
    echo "build status=$status"
    echo "build output=$output"
    [ "$status" -eq 0 ]

    # Test that bump will update version environment
    run wrench bump minor

    echo "bump status=$status"
    echo "bump output=$output"
    [ "$status" -eq 0 ]

    run wrench run current_version
    echo "run status=$status"
    echo "run lines[1]=${lines[1]}"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" == *"CURRENT VERSION v0.2.0"* ]]
}
