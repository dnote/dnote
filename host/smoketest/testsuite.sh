#!/usr/bin/env bash
# testsuite.sh runs the smoke tests for a self-hosted instance.
# It is meant to be run inside a virtual machine which has been
# set up by an entry script.
set -eu

echo 'Running a smoke test'

sudo -u postgres dropdb dnote
sudo -u postgres createdb dnote

cd /vagrant

tar -xvf dnote_server_integration_test_linux_amd64.tar.gz

GO_ENV=PRODUCTION \
  DBHost=localhost \
  DBPort=5432 \
  DBName=dnote \
  DBUser=postgres \
  DBPassword="" \
  WebURL=localhost:3000 \
  ./dnote-server -port 2300 start & sleep 3

assert_http_status() {
  url=$1
  expected=$2

  echo "======== [TEST CASE] asserting response status code for $url ========"

  got=$(curl --write-out %"{http_code}" --silent --output /dev/null "$url")

  if [ "$got" != "$expected" ]; then
    echo "======== ASSERTION FAILED ========"
    echo "status code for $url: expected: $expected got: $got"
    echo "=================================="
    exit 1
  fi
}

assert_http_status http://localhost:2300 "200"
assert_http_status http://localhost:2300/api/health "200"

echo "======== [SUCCESS] TEST PASSED! ========"
