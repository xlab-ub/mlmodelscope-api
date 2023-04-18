# MLModelScope API components

This repository provides the main parts of the MLModelScope API

## Running

To run a local version of the API you will need to have Docker and Docker
Compose installed. The following command will build the latest API image
and run it alongside the other necessary system components:

`docker compose up -d --build`

The additional components launched are:

* RabbitMQ - message queue providing communication between the API and ML agents
* PostgreSQL - the database
  * The database is initialized from the file `docker/data/c3sr-bootstrap.sql.gz`
* Companion - assists in cloud storage uploads, see below for details
* Traefik - reverse proxy, see below for details
* Consul - service discovery
* A suite of services to support monitoring with Prometheus/Grafana

## Deployment

The `DOCKER_REGISTRY` environment variable must be set to build or pull
the correct image tags for development, staging, or production. The `.env`
file sets this to `c3sr` by default so that images will be tagged and
pulled from the C3SR namespace on Docker Hub. Change this if you want to
use a private registry to host your own modified images.

This repository contains a Github workflow that will automatically build and
push an API image to the Github Container Registry each time new commits
are pushed to Github on the `master` branch.

You can read more about the Docker Compose configuration [here](docs/docker-compose.md).

## API

The `/api` directory contains an application that provides most of
the API endpoints for mlmodelscope.

### Running unit tests

To run the unit tests, change to the `/api` directory and run:

```bash
go test ./...
```

Add the `-v` flag to see detailed output from the tests:

```bash
go test -v ./...
```

### Running integration tests

To run the integration tests, change to the `/api` directory and run:

```bash
scripts/run-integration-tests.sh
```

This script will start the required services (RabbitMQ, PostgreSQL, and a Mock agent) in Docker containers and run the tests. When the tests are complete the containers will be stopped and removed.

### Debugging in a container

It is possible to debug the API endpoints while they run in a container
(this can be useful to test behavior when the API is running on a Docker
network alongside ML agents.) To enable debugging in the container, run
the API from the `docker/Dockerfile.api-debug` Dockerfile. This Dockerfile
creates a container that runs the API app with the [Delve](https://github.com/go-delve/delve) 
debugger attached. Delve listens on port 2345, which is exposed to the host
machine. The API itself will not begin running until a debugging client is
attached to Delve.

## Companion

[Companion](https://uppy.io/docs/companion/) is a service used to enable direct
uploads to a cloud storage provider. It requires additional environment variables
for configuration:

    * COMPANION_AWS_KEY
    * COMPANION_AWS_SECRET
    * COMPANION_AWS_BUCKET
    * COMPANION_AWS_REGION

In local development these variables should be provided in an environment
file named `.env.companion` in the project root folder. As this file will
likely contain private credentials it should **never** be committed to source
control!

## Traefik

[Traefik](https://doc.traefik.io/traefik/) is used as a reverse proxy for local
development to provide services at URLs such as http://api.local.mlmodelscope.org/.
If you are running a local copy of the
[MLModelScope React App](https://github.com/c3sr/mlmodelscope) on port 3000, Traefik
will proxy that at http://local.mlmodelscope.org/.

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