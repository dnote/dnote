#!/usr/bin/env bash
set -eux

sudo apt-get update
sudo apt-get install -y htop git wget build-essential inotify-tools

# Install Chrome
wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo apt-key add
echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' | sudo tee /etc/apt/sources.list.d/google-chrome.list
sudo apt-get -y update
sudo apt-get install -y google-chrome-stable
