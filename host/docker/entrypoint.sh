#!/bin/sh

# Set default DBPath to /data if not specified
export DBPath=${DBPath:-/data/dnote.db}

exec "$@"
