#!/usr/bin/env bash
# dump_schema.sh dumps the current system's dnote schema
set -eux

sqlite3 ~/.dnote/dnote.db .schema
