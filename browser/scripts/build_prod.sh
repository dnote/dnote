#!/bin/bash
# build_prod.sh builds distributable archive for the addon
# remember to bump version in package.json
set -eux

# clean
yarn clean

# chrome
yarn build:chrome
yarn package:chrome
# firefox
yarn build:firefox
yarn package:firefox
