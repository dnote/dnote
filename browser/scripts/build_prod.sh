#!/bin/bash
# build_prod.sh builds distributable archive for the addon
# remember to bump version in package.json
set -eux

# clean
npm run clean

# chrome
npm run build:chrome
npm run package:chrome
# firefox
npm run build:firefox
npm run package:firefox
