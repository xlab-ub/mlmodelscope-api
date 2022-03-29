#!/usr/bin/env bash

cd `git rev-parse --show-toplevel`
pg_dump -h localhost -p 15432 -U c3sr -f ./docker/data/c3sr-bootstrap.sql
gzip ./docker/data/c3sr-bootstrap.sql