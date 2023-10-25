#!/usr/bin/env bats

setup () {
    BATS_TMP_DIR=$(mktemp -d .wrench-bats.XXXXX)

    cp -r "$BATS_TEST_DIRNAME/../examples/simple" "$BATS_TMP_DIR/simple"

    pushd $BATS_TMP_DIR
    cd simple

    git init
    git config user.name "Your Name"
    git config user.email "you@example.com"
    git add .
    git commit -m "Initial commit"
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
    echo "bump output=$output"
    echo "bump status=$status"
    [ "$status" -eq 0 ]
    [ "$output" = "Released v0.1.0" ]

    # verify tag exists
    run git rev-list v0.1.0..
    echo "tag output=$output"
    echo "tag status=$status"
    [ "$status" -eq 0 ]

    # verify image exists
    run docker inspect example/simple:v0.1.0
    echo "image output=$output"
    echo "image status=$status"
    [ "$status" -eq 0 ]
}

@test "BUMP: missing docker snapshot image" {
    run wrench bump minor
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 1 ]
    [[ "$output" =~ Docker\ image\ for\ revision\ [^\ ]*\ could\ not\ be\ found ]]
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

    # verify tag exists
    run git rev-list v2.0.0..
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

    # verify tag exists
    run  git rev-list v1.3.0..
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

    # verify tag exists
    run git rev-list v1.2.34..
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]

    # verify image exists
    run docker inspect example/simple:v1.2.34
    echo "output=$output"
    echo "status=$status"
    [ "$status" -eq 0 ]
}

@test "BUMP: Chained changes from root commit" {
    # Create chain of changes and save their gitsha
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit A" --allow-empty
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit B" --allow-empty
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit C" --allow-empty
    # COMMIT_C=$(git rev-parse --short HEAD)
    run wrench build
    [ "$status" -eq 0 ]


    VERSION=0
    # Loop through all short sha listed from initial commit
    while read -r line; do
        # Increment to keep track of what should be current version
        VERSION=$((VERSION+1))

        # Checkout current change
        git checkout $line

        # Test bump
        run wrench bump minor
        echo "bump output=$output"
        echo "bump status=$status"
        [ "$status" -eq 0 ]
        [ "$output" = "Released v0.$VERSION.0" ]

        # verify tag exists
        run git rev-list v0.$VERSION.0..
        echo "tag output=$output"
        echo "tag status=$status"
        [ "$status" -eq 0 ]

        # verify image exists
        run docker inspect example/simple:v0.$VERSION.0
        echo "image output=$output"
        echo "image status=$status"
        [ "$status" -eq 0 ]

    done <<< "$(git log --pretty=format:'%h' --reverse)"
}

@test "BUMP: Chained changes from tagged root commit" {
    # Create chain of changes and save their gitsha
    git tag -a v0.1.0 -m "Release v0.1.0"
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit A" --allow-empty
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit B" --allow-empty
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit C" --allow-empty
    # COMMIT_C=$(git rev-parse --short HEAD)
    run wrench build
    [ "$status" -eq 0 ]


    VERSION=1
    # Loop through all short sha listed from initial commit
    while read -r line; do
        # Increment to keep track of what should be current version
        VERSION=$((VERSION+1))

        # Checkout current change
        git checkout $line

        # Test bump
        run wrench bump minor
        echo "bump output=$output"
        echo "bump status=$status"
        [ "$status" -eq 0 ]
        [ "$output" = "Released v0.$VERSION.0" ]

        # verify tag exists
        run git rev-list v0.$VERSION.0..
        echo "tag output=$output"
        echo "tag status=$status"
        [ "$status" -eq 0 ]

        # verify image exists
        run docker inspect example/simple:v0.$VERSION.0
        echo "image output=$output"
        echo "image status=$status"
        [ "$status" -eq 0 ]

    done <<< "$(git log --pretty=format:'%h' --reverse | tail -n +2)"
}

@test "BUMP: Chained changes from tagged patch release" {
    # Create chain of changes and save their gitsha
    git commit -m "Some patch release" --allow-empty
    git tag -a v1.10.4 -m "Patch release v1.10.4"

    git commit -m "Commit A" --allow-empty
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit B" --allow-empty
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit C" --allow-empty
    # COMMIT_C=$(git rev-parse --short HEAD)
    run wrench build
    [ "$status" -eq 0 ]


    VERSION=10
    # Loop through all short sha listed from initial commit
    while read -r line; do
        # Increment to keep track of what should be current version
        VERSION=$((VERSION+1))

        # Checkout current change
        git checkout $line

        # Test bump
        run wrench bump minor
        echo "bump output=$output"
        echo "bump status=$status"
        [ "$status" -eq 0 ]
        [ "$output" = "Released v1.$VERSION.0" ]

        # verify tag exists
        run git rev-list v1.$VERSION.0..
        echo "tag output=$output"
        echo "tag status=$status"
        [ "$status" -eq 0 ]

        # verify image exists
        run docker inspect example/simple:v1.$VERSION.0
        echo "image output=$output"
        echo "image status=$status"
        [ "$status" -eq 0 ]

    done <<< "$(git log --pretty=format:'%h' --reverse | tail -n +3)"
}

@test "BUMP: Chained changes from root commit with old history" {
    # Create some old history
    for i in `seq 1 7`; do
        git commit -m "Commit $i" --allow-empty
    done

    # Create chain of changes and save their gitsha
    git commit -m "Commit A" --allow-empty
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit B" --allow-empty
    run wrench build
    [ "$status" -eq 0 ]

    git commit -m "Commit C" --allow-empty
    # COMMIT_C=$(git rev-parse --short HEAD)
    run wrench build
    [ "$status" -eq 0 ]


    VERSION=0
    # Loop through our 3 changes and test bump
    while read -r line; do
        # Increment to keep track of what should be current version
        VERSION=$((VERSION+1))

        # Checkout current change
        git checkout $line

        # Test bump
        run wrench bump minor
        echo "bump output=$output"
        echo "bump status=$status"
        [ "$status" -eq 0 ]
        [ "$output" = "Released v0.$VERSION.0" ]

        # verify tag exists
        run git rev-list v0.$VERSION.0..
        echo "tag output=$output"
        echo "tag status=$status"
        [ "$status" -eq 0 ]

        # verify image exists
        run docker inspect example/simple:v0.$VERSION.0
        echo "image output=$output"
        echo "image status=$status"
        [ "$status" -eq 0 ]

    done <<< "$(git log --pretty=format:'%h' --reverse | tail -n 3)"
}
