#!/usr/bin/env bash

/docker-entrypoint.sh postgres &

export DBHost=localhost
export DBPort=5432
export DBName="$POSTGRES_DB"
export DBUser="$POSTGRES_USER"
export DBPassword="$POSTGRES_PASSWORD"
export WebURL=localhost:3000

until PGPASSWORD="$POSTGRES_PASSWORD" psql -h localhost -U "$POSTGRES_USER" "$DBName" &> /dev/null; do
  echo "Waiting for Postgres server to start..."
  sleep 3
done

exec ./dnote-server start
