#!/usr/bin/env bash
# build.sh builds styles
set -ex

dir=$(dirname "${BASH_SOURCE[0]}")
serverDir="$dir/../.."
outputDir="$serverDir/static"
inputDir="$dir/src"

rm -rf "${outputDir:?}/*"

"$dir/../node_modules/.bin/sass" --version

task="$dir/../node_modules/.bin/sass \
  --style compressed \
  --source-map \
  $inputDir:$outputDir"

# compile first then watch
eval "$task"

if [[ "$1" == "true" ]]; then
  eval "$task --watch --poll"
fi
