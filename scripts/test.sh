#!/bin/bash

# run_server_test.sh runs server test files sequentially
# https://stackoverflow.com/questions/23715302/go-how-to-run-tests-for-multiple-packages

# clear tmp dir in case not properly torn down
rm -rf ./tmp

# run test
go test ./... \
  -p 1\
  --tags "fts5"
