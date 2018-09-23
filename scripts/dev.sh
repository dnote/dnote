#!/bin/bash

# dev.sh builds a new binary and replaces the old one in the PATH with it

rm "$(which dnote)" $GOPATH/bin/cli
go install -ldflags "-X main.apiEndpoint=http://127.0.0.1:5000" --tags "darwin" .
ln -s $GOPATH/bin/cli /usr/local/bin/dnote
