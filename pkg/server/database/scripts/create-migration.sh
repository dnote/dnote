#!/usr/bin/env bash
# create-migration.sh creates a new SQL migration file for the
# server side Postgres database using the sql-migrate tool.
set -eux

is_command () {
  command -v "$1" >/dev/null 2>&1;
}

if ! is_command sql-migrate; then
  echo "sql-migrate is not found. Please run install-sql-migrate.sh"
  exit 1
fi

if [ "$#" == 0 ]; then
  echo "filename not provided"
  exit 1
fi

filename=$1
sql-migrate new -config=./sql-migrate.yml "$filename"
