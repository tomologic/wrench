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
}

teardown () {
    popd
    rm -rf $BATS_TMP_DIR
}

@test "CONFIG: detect project name" {
    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    expected="Name: $(basename $PWD)"

    echo "expected=$expected"
    echo "out=$out"
    echo $out | grep "$expected"
}

@test "CONFIG: detect project organization" {
    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    num_punctation="$(hostname -f | grep -o '\.' | wc -l | tr -d '[[:space:]]')"

    if [ $num_punctation == "1" ]; then
        # expect host.local => local
        domain="$(hostname -f | cut -d. -f2)"
        expected="Organization: $domain"
    else
        # example.com          => example
        # host.example.com     => example
        # host.2.example.com   => example
        # host.2...example.com => example
        domain="$(hostname -f | rev | cut -d. -f2 | rev)"
        expected="Organization: $domain"
    fi

    echo "hostname=$(hostname -f)"
    echo "num_punctation=$num_punctation"
    echo "expected=$expected"
    echo "out=$out"
    echo $out | grep "$expected"
}

@test "CONFIG: detect project initial version" {
    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    expected="Version: v0.0.0-1-g$(git rev-parse --short HEAD)"

    echo "expected=$expected"
    echo "out=$out"
    echo $out | grep "$expected"
}

@test "CONFIG: detect project release version" {
    git tag -a v0.10.0 -m "release v0.10.0"

    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    expected="Version: v0.10.0"

    echo "expected=$expected"
    echo "out=$out"
    echo $out | grep "$expected"
}

@test "CONFIG: detect project snapshot version" {
    git tag -a v0.10.0 -m "release v0.10.0"
    touch b
    git add b
    git commit -m "add b"

    ret=0
    out=$(wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 0 ]

    expected="Version: v0.10.0-1-g$(git rev-parse --short HEAD)"

    echo "expected=$expected"
    echo "out=$out"
    echo $out | grep "$expected"
}
