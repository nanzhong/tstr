#!/usr/bin/env bash

buf generate
buf generate buf.build/envoyproxy/protoc-gen-validate
