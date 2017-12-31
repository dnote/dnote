#!/bin/bash

rm $(which dnote) && go install -ldflags "-X main.apiEndpoint=http://127.0.0.1:5000" . && ln -s $GOPATH/bin/cli /usr/local/bin/dnote
