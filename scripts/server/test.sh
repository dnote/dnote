#!/usr/bin/env bash
# test.sh runs server tests. It is to be invoked by other scripts that set
# appropriate env vars.
set -ex

dir=$(realpath "$(dirname "${BASH_SOURCE[0]}")")
pushd "$dir/../../pkg/server"

function run_test {
  if [ -z "$1" ]; then
    go test -tags "fts5" ./... -cover
  else
    go test -tags "fts5" -run "$1" -cover
  fi
}

if [ "${WATCH-false}" == true ]; then
  set +e
  while inotifywait --exclude .swp -e modify -r .; do run_test; done;
  set -e
else
  run_test "$1"
fi

popd
