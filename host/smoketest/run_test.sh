#!/usr/bin/env bash
# run_test.sh builds a fresh server image, and mounts it on a fresh
# virtual machine and runs a smoke test.
set -ux

dir=$(dirname "${BASH_SOURCE[0]}")
projectDir="$dir/../.."

# build
pushd "$projectDir"
make version=integration_test build-server
popd

pushd "$dir"
volume="$dir/volume"
rm -rf "$volume"
mkdir -p "$volume"
cp "$projectDir/build/server/dnote_server_integration_test_linux_amd64.tar.gz" "$volume"
cp "$dir/testsuite.sh" "$volume"

# test
vagrant up

if ! vagrant ssh -c "/vagrant/testsuite.sh"; then
  echo "Test failed. Please see the output."
  exit 1
fi

vagrant halt
popd

