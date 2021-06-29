#!/usr/bin/env bash

cd `git rev-parse --show-toplevel`
source "scripts/functions.sh"

USAGE="Usage: ${0} <target>

    target        The service to build into a container, i.e. 'uploader':
                      > build.sh uploader"

TARGET="$1"
TARGET_DIR="./$TARGET"
DOCKERFILE="./docker/Dockerfile.$1"


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