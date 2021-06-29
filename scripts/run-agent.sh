#!/usr/bin/env sh

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
    --network=host \
    -P \
    -v "`pwd`/.carml_config.yml:/root/.carml_config.yml" \
    -v "/tmp/results:/go/src/github.com/rai-project/mlmodelscope/results" \
    c3sr/$AGENT:amd64-cpu-latest serve -l -d -v

cleanup_and_exit