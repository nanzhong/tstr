#!/usr/bin/env bash

set -euo pipefail

export DATABASE_URL="postgres://tstr@127.0.0.1:5432/tstr_development?sslmode=disable"
export TEST_DATABASE_URL="postgres://tstr@127.0.0.1:5432/tstr_test?sslmode=disable"

export dev_state_path=".dev"
export dev_overmind_sock="$dev_state_path/overmind.sock"
