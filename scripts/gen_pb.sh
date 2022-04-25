#!/usr/bin/env bash

shopt -s globstar

protoc \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    **/*.proto
