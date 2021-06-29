# mlmodelscope API components

This repository provides the main parts of the mlmodelscope API

## API

The `/api` directory contains an application that provides most of
the API endpoints for mlmodelscope.

### Debugging in a container

To debug the API endpoints while they run in a container (this can
be useful to test behavior when the API is running on a Docker
network alongside ML agents), you can build a Docker container using
`docker/Dockerfile.api-debug`.

This Dockerfile creates a container that runs the API app with the
[Delve](https://github.com/go-delve/delve) debugger attached. Delve
listens on port 2345, which you should expose to your host machine.
The API itself will not begin running until a debugging client is
attached to Delve.

The `docker-compose.yml` can be updated to use the debug Dockerfile
to build the API container by changing the **build** section of the
**api** service to use the correct Dockerfile, and exposing port 2345
in the **ports** section.

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