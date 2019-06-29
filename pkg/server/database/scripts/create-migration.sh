#!/bin/bash
# create-migration.sh creates a new SQL migration file for the
# server side Postgres database using the sql-migrate tool.
set -eux

if [ "$#" == 0 ]; then
  echo "filename not provided"
  exit 1
fi

filename=$1
sql-migrate new -config=migrate.yml "$filename"
