#!/usr/bin/env bash

dbmate --url 'postgres://tstr@127.0.0.1:5432/tstr_development?sslmode=disable' up

sqlc generate
