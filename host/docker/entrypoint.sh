#!/bin/sh

wait_for_db() {
  HOST=${DBHost:-postgres}
  PORT=${DBPort:-5432}
  echo "Waiting for the database connection..."

  attempts=0
  max_attempts=10
  while [ $attempts -lt $max_attempts ]; do
    nc -z "${HOST}" "${PORT}" 2>/dev/null && break
    echo "Waiting for db at ${HOST}:${PORT}..."
    sleep 5
    attempts=$((attempts+1))
  done

  if [ $attempts -eq $max_attempts ]; then
    echo "Timed out while waiting for db at ${HOST}:${PORT}"
    exit 1
  fi
}

wait_for_db

exec "$@"
