#!/usr/bin/env sh

cd `git rev-parse --show-toplevel`/api

source ./scripts/functions.sh

start_integration_containers
