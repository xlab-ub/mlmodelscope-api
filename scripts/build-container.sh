#!/usr/bin/env bash

cleanup_and_exit () {
    cd - 2&>1 /dev/null
    exit $1
}

die_with_message () {
    echo "$1"
    cleanup_and_exit 1
}

USAGE="Usage: ${0} <target>

    target        The service to build into a container, i.e. 'uploader':
                      > build.sh uploader"

TARGET="./$1"
DOCKERFILE="./docker/Dockerfile.$1"

cd `git rev-parse --show-toplevel`

if [ -z $TARGET ]; then
    die_with_message "$USAGE"
fi

if [ ! -d "$TARGET" ]; then
    die_with_message "Target ${TARGET} does not exist"
fi

if [ ! -f "$DOCKERFILE" ]; then
    die_with_message "Dockerfile ${DOCKERFILE} does not exist"
fi

docker build -t "mlmodelscope-api/$1:latest" --file "$DOCKERFILE" "$TARGET"

cleanup_and_exit