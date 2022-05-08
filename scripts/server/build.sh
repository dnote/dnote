#!/usr/bin/env bash
set -eux

dir=$(dirname "${BASH_SOURCE[0]}")

version=$1
projectDir="$dir/../.."
basedir="$projectDir/pkg/server"
outputDir="$projectDir/build/server"

command_exists () {
  command -v "$1" >/dev/null 2>&1;
}

if ! command_exists shasum; then
  echo "please install shasum"
  exit 1
fi
if [ $# -eq 0 ]; then
  echo "no version specified."
  exit 1
fi
if [[ $1 == v* ]]; then
  echo "do not prefix version with v"
  exit 1
fi

build() {
  platform=$1
  arch=$2

  pushd "$basedir"

  destDir="$outputDir/$platform-$arch"
  mkdir -p "$destDir"

  # build binary
  moduleName="github.com/dnote/dnote"
  ldflags="-X '$moduleName/pkg/server/buildinfo.CSSFiles=main.css' -X '$moduleName/pkg/server/buildinfo.JSFiles=main.js' -X '$moduleName/pkg/server/buildinfo.Version=$version' -X '$moduleName/pkg/server/buildinfo.Standalone=true'"

  GOOS="$platform" \
  GOARCH="$arch" go build \
    -o "$destDir/dnote-server" \
    -ldflags "$ldflags" \
    "$basedir"/*.go

  popd

  # build tarball
  tarballName="dnote_server_${version}_${platform}_${arch}.tar.gz"
  tarballPath="$outputDir/$tarballName"

  cp "$projectDir/licenses/AGPLv3.txt" "$destDir"
  cp "$basedir/README.md" "$destDir"
  tar -C "$destDir" -zcvf "$tarballPath" "."
  rm -rf "$destDir"

  # calculate checksum
  pushd "$outputDir"
  shasum -a 256 "$tarballName" >> "$outputDir/dnote_${version}_checksums.txt"
  popd

}

build linux amd64
build linux arm64
