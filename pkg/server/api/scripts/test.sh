#!/usr/bin/env bash
# test.sh runs api tests. It is to be invoked by other scripts that set
# appropriate env vars.
set -eux

pushd "$GOPATH"/src/github.com/dnote/dnote/pkg/server/api

if [ "${WATCH-false}" == true ]; then
  set +e
  while inotifywait --exclude .swp -e modify -r .; do go test ./... -cover -p 1; done;
  set -e
else
  go test ./... -cover -p 1
fi

popd
