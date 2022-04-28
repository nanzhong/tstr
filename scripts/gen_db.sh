#!/usr/bin/env bash

dbmate --url 'postgres://tstr@127.0.0.1:5432/tstr_development?sslmode=disable' up

pggen gen go \
      --postgres-connection "user=tstr host=127.0.0.1 dbname=tstr_development" \
      --query-glob db/queries.sql \
      --go-type 'varchar=string' \
      --go-type '_varchar=[]string' \
      --go-type 'int4=int' \
      --go-type 'timestamptz=time.Time' \
      --go-type 'uuid=string' \
