#!/usr/bin/env bats

setup () {
    # Get base image for bogus docker image
    BASE_IMAGE="busybox"
    if docker history -q $BASE_IMAGE; then
        echo "$BASE_IMAGE already exists"
    else
        docker pull $BASE_IMAGE
    fi

    TEST_IMAGE_NAME="bats/wrenchtests"
    TEST_IMAGE_TAG="v1.0.0"
    TEST_IMAGE="$TEST_IMAGE_NAME:$TEST_IMAGE_TAG"
    # Create bogus docker image for tests
    if docker history -q $TEST_IMAGE > /dev/null 2>&1; then
        echo "$TEST_IMAGE already exists"
    else
        echo "Fetching $TEST_IMAGE"
        docker tag $BASE_IMAGE $TEST_IMAGE
    fi

    # Pull registry image if missing
    docker history -q registry:2 > /dev/null 2>&1 || {
        docker pull registry:2
    }

    # Start a docker registry
    REGISTRY_ID=$(docker run -P -d registry:2)

    # Get port for registry
    REGISTRY_PORT=$(docker inspect \
        -f '{{index .NetworkSettings.Ports "5000/tcp" 0 "HostPort"}}' \
        $REGISTRY_ID)

    REGISTRY="127.0.0.1:$REGISTRY_PORT"

    # Check if docker machine is in use
    if [ -n "$DOCKER_MACHINE_NAME" ]; then
        # Get host for registry
        DOCKER_MACHINE_IP=$(docker-machine ip $DOCKER_MACHINE_NAME)

        REGISTRY_API_URL="$DOCKER_MACHINE_IP:$REGISTRY_PORT"
    else
        # Assume docker on 127.0.0.1
        REGISTRY_API_URL="127.0.0.1:$REGISTRY_PORT"
    fi
    echo "REGISTRY=$REGISTRY"
    echo "REGISTRY_API_URL=$REGISTRY_API_URL"

    echo "# SETUP DONE"
    echo
}

teardown () {
    echo
    echo "# TEARDOWN START"
    # Stop docker registry
    docker rm -f $REGISTRY_ID

    # Remove all images assosiated with temporary registries
    # TODO: remove and add verification that wrench do the cleanup
    test_images=$(docker images | grep -o "$REGISTRY.*bats\/\S*[[:space:]]*\S*" | tr -s ' ' | sed 's/\ /:/')
    if [ -n "$test_images" ]; then
        docker rmi $test_images
    fi
}

json_stemming() {
    # Json diff is hard, lets not do it. json_stemming will do following steps:
    # 1. Replace all {},:" with spaces
    #   '{"hello":"world","foo/bar":"v1.0.0"} -> '  hello   world   foo/bar   v1.0.0  '
    # 2. Xargs output each word on independent row
    # 3. Sort the result
    # We now have a string that can be compared that should be the same with
    # the same json structure
    echo $1 | sed 's/[{},:"]/ /g' | xargs -n1 | sort
}

@test "PUSH: Push image to registry" {
    # TODO: Replace with wrench command
    # wrench push
    docker tag $TEST_IMAGE $REGISTRY/$TEST_IMAGE
    docker push $REGISTRY/$TEST_IMAGE

    actual=$(json_stemming $(curl "$REGISTRY_API_URL/v2/$TEST_IMAGE_NAME/tags/list"))
    echo "actual=$actual"

    expected=$(json_stemming '{"name":"bats/wrenchtests","tags":["v1.0.0"]}')
    echo "expected=$expected"

    [ "$actual" == "$expected" ]
}

@test "PUSH: Additional tag" {
    # TODO: Replace with wrench command
    # wrench push --additional-tags additional
    docker tag $TEST_IMAGE $REGISTRY/$TEST_IMAGE
    docker tag $TEST_IMAGE "$REGISTRY/$TEST_IMAGE_NAME:additional"
    docker push $REGISTRY/$TEST_IMAGE
    docker push "$REGISTRY/$TEST_IMAGE_NAME:additional"

    actual=$(json_stemming $(curl "$REGISTRY_API_URL/v2/$TEST_IMAGE_NAME/tags/list"))
    echo "actual=$actual"

    expected=$(json_stemming '{"name":"bats/wrenchtests","tags":["v1.0.0","additional"]}')
    echo "expected=$expected"

    [ "$actual" == "$expected" ]
}

@test "PUSH: Several additional tags" {
    # TODO: Replace with wrench command
    # wrench push --additional-tags additional,latest,tag3
    docker tag $TEST_IMAGE $REGISTRY/$TEST_IMAGE
    docker tag $TEST_IMAGE "$REGISTRY/$TEST_IMAGE_NAME:additional"
    docker tag $TEST_IMAGE "$REGISTRY/$TEST_IMAGE_NAME:latest"
    docker tag $TEST_IMAGE "$REGISTRY/$TEST_IMAGE_NAME:tag3"
    docker push $REGISTRY/$TEST_IMAGE
    docker push "$REGISTRY/$TEST_IMAGE_NAME:additional"
    docker push "$REGISTRY/$TEST_IMAGE_NAME:latest"
    docker push "$REGISTRY/$TEST_IMAGE_NAME:tag3"

    actual=$(json_stemming $(curl "$REGISTRY_API_URL/v2/$TEST_IMAGE_NAME/tags/list"))
    echo "actual=$actual"

    expected=$(json_stemming '{"name":"bats/wrenchtests","tags":["v1.0.0","latest","additional","tag3"]}')
    echo "expected=$expected"

    [ "$actual" == "$expected" ]
}
