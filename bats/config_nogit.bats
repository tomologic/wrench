#!/usr/bin/env bats

setup () {
    BATS_TMP_DIR=$(mktemp -d .wrench-bats.XXXXX)
    pushd $BATS_TMP_DIR
}

teardown () {
    popd
    rm -rf $BATS_TMP_DIR
}

@test "CONFIG: git not installed" {
    # Fake 127 exit code
    echo "#!/bin/bash" > git
    echo "exit 127" >> git
    chmod +x git

    ret=0
    out=$(PATH=$PWD:$PATH wrench config) || ret=$?

    echo "ret=$ret"
    [ "$ret" -eq 1 ]

    echo "out=$out"
    echo $out | grep "No git executable found"
}
