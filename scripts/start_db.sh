#!/usr/bin/env bash

if [ ! -d "$(pwd)/.dev/db" ]; then
    initdb -U tstr "$(pwd)/.dev/db/data"
fi

postgres -D "$(pwd)/.dev/db/data" -k "$(pwd)/.dev/db"
