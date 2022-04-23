#!/usr/bin/env bash
# shellcheck disable=SC1090
# dev.sh builds and starts development environment
set -eux -o pipefail

dir=$(dirname "${BASH_SOURCE[0]}")
basePath="$dir/../.."
serverPath="$basePath/pkg/server"

# load env
set -a
dotenvPath="$serverPath/.env.dev"
source "$dotenvPath"
set +a

# copy assets
mkdir -p "$basePath/pkg/server/static"
cp "$basePath"/pkg/server/assets/static/* "$basePath/pkg/server/static"
# run asset pipeline in the background
(cd "$basePath/pkg/server/assets/" && "$basePath/pkg/server/assets/styles/build.sh" true ) &
(cd "$basePath/pkg/server/assets/" && "$basePath/pkg/server/assets/js/build.sh") &

# run server
moduleName="github.com/dnote/dnote"
ldflags="-X '$moduleName/pkg/server/buildinfo.CSSFiles=main.css' -X '$moduleName/pkg/server/buildinfo.JSFiles=main.js' -X '$moduleName/pkg/server/buildinfo.Version=dev'"
task="go run -ldflags \"$ldflags\" main.go start -port 3000"

(
  cd "$basePath/pkg/watcher" && \
  go run main.go \
  --task="$task" \
  --context="$serverPath" \
  "$serverPath"
)
