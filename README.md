# mlmodelscope API components

This repository provides the main parts of the mlmodelscope API

## Deployment

The `DOCKER_REGISTRY` environment variable must be set to build or pull
the correct image tags for development, staging, or production. The `.env`
file sets this to `c3sr` by default so that images will be tagged and
pulled from the C3SR namespace on Docker Hub. Change this if you want to
use a private registry to host your own modified images.

## API

The `/api` directory contains an application that provides most of
the API endpoints for mlmodelscope.

### Debugging in a container

It is possible to debug the API endpoints while they run in a container
(this can be useful to test behavior when the API is running on a Docker
network alongside ML agents.) To enable debugging in the container, run
the API from the `docker/Dockerfile.api-debug` Dockerfile. This Dockerfile
creates a container that runs the API app with the [Delve](https://github.com/go-delve/delve) 
debugger attached. Delve listens on port 2345, which is exposed to the host
machine. The API itself will not begin running until a debugging client is
attached to Delve.

## Uploader

The `/uploader` directory contains an application that provides a file
upload endpoint backed by [tusd](https://github.com/tus/tusd).

## Running an agent alongside the API

The `scripts/run-agent.sh` script will run an agent container for one of the
following ML frameworks:

    * mxnet
    * onnxruntime
    * pytorch
    * tensorflow

The `docker/carml-config-examle.yml` file will be copied to `.carml_config.yml` and
that file will be mapped into the running container as a Docker volume. If you
need to modify the configuration in any way, you should edit the `.carml_config.yml`
file and **not** `docker/carml-config-example.yml`.