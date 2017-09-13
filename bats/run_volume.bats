#!/usr/bin/env bats

setup () {
    BATS_TMP_DIR=$(mktemp -d .wrench-bats.XXXXX)

    cp -r "$BATS_TEST_DIRNAME/../examples/simple" "$BATS_TMP_DIR/origin"

    pushd $BATS_TMP_DIR
    pushd origin

    # Add Run target which reads mounted file
    echo "  read_file:
    Cmd: cat /mounted_file
    Volumes:
      - $(pwd)/text_file:/mounted_file" >> wrench.yml

    # Add text file
    echo "confirmed" > text_file

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

@test "VOLUMES: Read from mounted volume in run target" {
    run wrench build
    echo "build status=$status"
    echo "build output=$output"
    [ "$status" -eq 0 ]

    run wrench run read_file
    echo "run status=$status"
    echo "run lines[1]=${lines[1]}"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" == "confirmed"* ]]
}
