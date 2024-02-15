#!/usr/bin/env bash
set -eux

sudo apt-get update
sudo apt-get install -y htop git wget build-essential inotify-tools

# Install Chrome
wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo apt-key add
echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' | sudo tee /etc/apt/sources.list.d/google-chrome.list
sudo apt-get -y update
sudo apt-get install -y google-chrome-stable

# Install dart-sass
dart_version=1.34.1
dart_tarball="dart-sass-$dart_version-linux-x64.tar.gz"
wget -q "https://github.com/sass/dart-sass/releases/download/$dart_version/$dart_tarball"
tar -xvzf "$dart_tarball" -C /tmp/
sudo install /tmp/dart-sass/sass /usr/bin
