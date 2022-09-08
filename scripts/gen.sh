#!/usr/bin/env bash

set -eo pipefail

scripts_path="${BASH_SOURCE%/*}"

# shellcheck source=scripts/common.sh
source "$scripts_path/common.sh"

options=$(getopt -l "gotags:" -o "t:" -a -- "$@")
eval set -- "$options"

gotags=""
while true
do
  case $1 in
    -t|--gotags)
      shift
      gotags=$1
      ;;
    --)
      shift
      break;;
  esac
  shift
done

run "Generating protobufs..." "$scripts_path/gen_pb.sh"
run "Generating db implementation..." "$scripts_path/gen_db.sh"
run "Running go generate..." go generate -tags="$gotags" ./...
