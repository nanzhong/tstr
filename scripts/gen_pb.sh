#!/usr/bin/env bash

shopt -s globstar

protoc \
    --proto_path=. \
    --proto_path=vendor/github.com/envoyproxy/protoc-gen-validate \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --validate_out="lang=go:." \
    --validate_opt="paths=source_relative" \
    api/**/*.proto
