#!/usr/bin/env bash

scripts_path="${BASH_SOURCE%/*}"

# run starts running a named command making sure to format stdout and stderr of
# the command with descriptive line prefixes.
# arguments: msg, command, [arg1 arg2 ...]
function run() {
  echo "┬─> $1"
  "${@:2}" > >(prepend "│o: ") 2> >(prepend "│e: " >&2)
  echo "└─> Done"
}

# prepend reads line on stdin and prepends a prefix to each line.
# arguments: prefix
function prepend() {
    while read -r line; do
        if [[ -n "$line" ]]; then
            echo "$1$line"
        fi
    done
}

run "Generating protobufs..." "$scripts_path/gen_pb.sh"
run "Generating db implementation..." "$scripts_path/gen_db.sh"
run "Running go generate..." go generate ./...
