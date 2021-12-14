#!/usr/bin/env sh

cd `git rev-parse --show-toplevel`/api

export MQ_HOST=localhost
export MQ_PORT=5672
export MQ_USER=user
export MQ_PASSWORD=password

printf "Starting RabbitMQ container"
docker run -d --rm --name mlmodelscope-api_rabbitmq-integration -p 5672:5672 --env RABBITMQ_DEFAULT_USER="${MQ_USER}" --env RABBITMQ_DEFAULT_PASS="${MQ_PASSWORD}" rabbitmq:3-alpine > /dev/null

while [ -z "`docker logs mlmodelscope-api_rabbitmq-integration | grep "started TCP listener"`" ]
do
    printf "."
    sleep 1
done

echo "\nRunning integration tests..."
go clean -testcache && go test -v -coverprofile integration.out --tags integration ./...

echo "Cleaning up RabbitMQ container..."
docker stop mlmodelscope-api_rabbitmq-integration > /dev/null

cd - 2>&1 >> /dev/null
