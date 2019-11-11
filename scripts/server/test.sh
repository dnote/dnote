#!/usr/bin/env bash
# test.sh runs server tests. It is to be invoked by other scripts that set
# appropriate env vars.
set -eux

dir=$(dirname $"{BASH_SOURCE[0]}")
pushd "$dir/../../pkg/server"

if [ "${WATCH-false}" == true ]; then
  set +e
  while inotifywait --exclude .swp -e modify -r .; do go test ./... -cover -p 1; done;
  set -e
else
  go test ./... -cover -p 1
fi

popd
