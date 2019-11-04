#!/usr/bin/env bash
set -ex

echo "export DNOTE=/go/src/github.com/dnote/dnote" >> /home/vagrant/.bash_profile
echo "cd /go/src/github.com/dnote/dnote" >> /home/vagrant/.bash_profile

# install dependencies
(cd /go/src/github.com/dnote/dnote && make install)

# set up database
sudo -u postgres createdb dnote
sudo -u postgres createdb dnote_test
# allow connection from host and allow to connect without password
sudo sed -i  "/port*/a listen_addresses = '*'" /etc/postgresql/11/main/postgresql.conf
sudo sed -i 's/host.*all.*.all.*md5/# &/' /etc/postgresql/11/main/pg_hba.conf
sudo sed -i "$ a host all all all trust" /etc/postgresql/11/main/pg_hba.conf
sudo service postgresql restart
