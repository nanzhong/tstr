#!/usr/bin/env bash

if [ ! -d "$(pwd)/.db" ]; then
    initdb -U tstr "$(pwd)/.db/data"
fi

postgres -D "$(pwd)/.db/data" -k "$(pwd)/.db"
