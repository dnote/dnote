#!/usr/bin/env bash
# build.sh builds styles
set -ex

dir=$(dirname "${BASH_SOURCE[0]}")
basePath="$dir/../../.."
serverDir="$dir/../.."
outputDir="$serverDir/static"
inputDir="$dir/src"

task="cp $inputDir/main.js $outputDir"


if [[ "$1" == "true" ]]; then
(
  cd "$basePath/watcher" && \
  go run main.go \
  --task="$task" \
  --context="$inputDir" \
  "$inputDir"
)
else
  eval "$task"
fi
