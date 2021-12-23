#!/usr/bin/env sh

cd `git rev-parse --show-toplevel`/api

export DB_DBNAME=integration
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_PASSWORD=password
export DB_PORT=5432
export DB_USER=user

export MQ_HOST=localhost
export MQ_PORT=5672
export MQ_USER=user
export MQ_PASSWORD=password

echo "Starting PostgreSQL container"
docker run \
    -d \
    --rm \
    --name mlmodelscope-api_postgres-integration \
    -p ${DB_PORT}:5432 \
    --env POSTGRES_USER="${DB_USER}" \
    --env POSTGRES_PASSWORD="${DB_PASSWORD}" \
    --env POSTGRES_DB="${DB_DBNAME}" \
    postgres:14.1 \
    > /dev/null

printf "Starting RabbitMQ container"
docker run \
    -d \
    --rm \
    --name mlmodelscope-api_rabbitmq-integration \
    -p ${MQ_PORT}:5672 \
    --env RABBITMQ_DEFAULT_USER="${MQ_USER}" \
    --env RABBITMQ_DEFAULT_PASS="${MQ_PASSWORD}" \
    rabbitmq:3-alpine \
    > /dev/null

while [ -z "`docker logs mlmodelscope-api_rabbitmq-integration | grep "started TCP listener"`" ]
do
    printf "."
    sleep 1
done
printf "\n"

echo "Starting Mock Agent container..."
docker run \
    -d \
    --rm \
    --name mlmodelscope-api_mock-agent-integration \
    --env MQ_HOST="host.docker.internal" \
    --env MQ_PORT="${MQ_PORT}" \
    --env MQ_USER="${MQ_USER}" \
    --env MQ_PASSWORD="${MQ_PASSWORD}" \
    registry.staging.mlmodelscope.org:5000/mock-agent:latest \
    > /dev/null

echo "\nRunning integration tests..."
go clean -testcache && go test -v -coverprofile integration.out --tags integration ./...

echo "Cleaning up PostgreSQL container..."
docker stop mlmodelscope-api_postgres-integration > /dev/null

echo "Cleaning up RabbitMQ container..."
docker stop mlmodelscope-api_rabbitmq-integration > /dev/null

echo "Cleaning up Mock Agent container..."
docker stop mlmodelscope-api_mock-agent-integration > /dev/null

cd - 2>&1 >> /dev/null
