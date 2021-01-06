#!/usr/bin/env bash
# test.sh runs server tests. It is to be invoked by other scripts that set
# appropriate env vars.
set -eux

dir=$(realpath "$(dirname "${BASH_SOURCE[0]}")")
pushd "$dir/../../pkg/server"

emailTemplateDir=$(realpath "$dir/../../pkg/server/mailer/templates/src")
export DNOTE_TEST_EMAIL_TEMPLATE_DIR="$emailTemplateDir"

function run_test {
  if [ -z "$1" ]; then
    go test ./... -cover -p 1
  else
    go test "$1" -cover -p 1
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
