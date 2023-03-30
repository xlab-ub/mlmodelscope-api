#!/usr/bin/env sh

# This script is used to run a single agent container for a given framework.
# It pulls the latest image from the GitHub Container Registry and runs it.
# You will need to be logged in to the registry to pull the image.
# See the following link for instructions to set up a Personal Access Token and use it to log in:
# https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry

cd `git rev-parse --show-toplevel`
source "scripts/functions.sh"

USAGE="Usage: $(basename $0) <framework>

    framework       The framework to run an agent for, one of:
                        * mxnet
                        * onnxruntime
                        * pytorch
                        * tensorflow"


TARGET=$(echo $1 | awk '/^(mxnet|onnxruntime|pytorch|tensorflow)$/{print $0}')
AGENT="$TARGET-agent"

if [ -z $TARGET ]; then
    die_with_message "$USAGE"
fi

if [ ! -f ".carml_config.yml" ]; then
    # copy the example config file to a real config file
    cp docker/carml-config-example.yml .carml_config.yml
fi

docker run \
    -d \
    --rm \
    --name="$AGENT" \
    --shm-size=1g \
    --ulimit memlock=-1 \
    --ulimit stack=67108864 \
    --privileged=true \
    --network=mlmodelscope-api_default \
    --env-file .env \
    -P \
    -v "`pwd`/.carml_config.yml:/root/.carml_config.yml" \
    -v "/tmp/results:/go/src/github.com/rai-project/mlmodelscope/results" \
    ghcr.io/c3sr/$AGENT:latest

cleanup_and_exit
