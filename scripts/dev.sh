#!/bin/bash

rm $(which dnote) && go install . && ln -s $GOPATH/bin/cli /usr/local/bin/dnote
