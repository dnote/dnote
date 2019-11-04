#!/usr/bin/env bash
# shellcheck disable=SC1091
set -eux

VERSION=1.13.4
OS=linux
ARCH=amd64

tarball=go$VERSION.$OS-$ARCH.tar.gz

wget -q https://dl.google.com/go/"$tarball"
sudo tar -C /usr/local -xzf "$tarball"
sudo tar -xf "$tarball"

sudo mkdir -p /go/src
sudo mkdir -p /go/bin
sudo mkdir -p /go/pkg
sudo chown -R vagrant:vagrant /go

GOPATH=/go
echo "export GOPATH=$GOPATH" >> /home/vagrant/.bash_profile
echo "export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin" >> /home/vagrant/.bash_profile
source /home/vagrant/.bash_profile

go version
go env
