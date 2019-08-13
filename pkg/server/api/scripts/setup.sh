#!/usr/bin/env bash
# setup.sh installs programs and depedencies necessary to run the project locally
# usage: ./setup.sh
set -eux

curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
go get github.com/githubnemo/CompileDaemon
dep ensure

