#!/usr/bin/env bash
set -ex

sudo apt-get install wget ca-certificates
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'

sudo apt-get update
sudo apt-get install -y postgresql-11

# set up database
sudo -u postgres createdb dnote
# allow connection from host and allow to connect without password
sudo sed -i  "/port*/a listen_addresses = '*'" /etc/postgresql/11/main/postgresql.conf
sudo sed -i 's/host.*all.*.all.*md5/# &/' /etc/postgresql/11/main/pg_hba.conf
sudo sed -i "$ a host all all all trust" /etc/postgresql/11/main/pg_hba.conf
sudo service postgresql restart
