#!/usr/bin/env bash
# run_test.sh builds a fresh server image, and mounts it on a fresh
# virtual machine and runs a smoke test. If a tarball path is not provided,
# this script builds a new version and uses it.
set -ex

# tarballPath is an absolute path to a release tarball containing the dnote server.
tarballPath=$1

dir=$(dirname "${BASH_SOURCE[0]}")
projectDir="$dir/../.."

# build
if [ -z "$tarballPath"  ]; then
  pushd "$projectDir"
  make version=integration_test build-server
  popd
  tarballPath="$projectDir/build/server/dnote_server_integration_test_linux_amd64.tar.gz"
fi

pushd "$dir"

# start a virtual machine
volume="$dir/volume"
rm -rf "$volume"
mkdir -p "$volume"
cp "$tarballPath" "$volume"
cp "$dir/testsuite.sh" "$volume"

vagrant up

# run tests
set +e
if ! vagrant ssh -c "/vagrant/testsuite.sh"; then
  echo "Test failed. Please see the output."
  vagrant halt
  exit 1
fi
set -e

vagrant halt
popd
