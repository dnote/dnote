#!/bin/bash
# dump_schema.sh dumps the current system's dnote schema to testutils package
# to be used while setting up tests

sqlite3 ~/.dnote/dnote.db .schema > ./testutils/fixtures/schema.sql
