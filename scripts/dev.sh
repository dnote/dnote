#!/bin/bash

# dev.sh builds a new binary and replaces the old one in the PATH with it

sudo rm "$(which dnote)" $GOPATH/bin/cli

# change --tags to "darwin" if on macOS
go install -ldflags "-X main.apiEndpoint=http://127.0.0.1:5000" --tags "linux" .

sudo ln -fs $GOPATH/bin/cli /usr/local/bin/dnote
