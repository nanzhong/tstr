#!/usr/bin/env bash

for table in access_tokens runners runs schema_migrations test_run_configs test_suites tests
do
psql -h 127.0.0.1 tstr_development -U tstr -c "copy (select row_to_json(t) from (select * from $table) t) to '$(pwd)/$table.json'"
done
