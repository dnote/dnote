#!/usr/bin/env bash
set -eux

PID=""

function cleanup {
  if [ "$PID" != "" ]; then
    kill "$PID"
  fi
}
trap cleanup EXIT

while true; do
  go build main.go
  ./main &
  PID=$!
  inotifywait -r -e modify .
  kill $PID
done


