#!/bin/bash

# dev.sh builds a new binary and replaces the old one in the PATH with it

rm "$(which dnote)" $GOPATH/bin/cli
make build-snapshot
ln -s $GOPATH/src/github.com/dnote/cli/dist/darwin_amd64/dnote /usr/local/bin/dnote
