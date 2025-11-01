#!/usr/bin/env bash
#
# build.sh compiles dnote binary for target platforms. It is resonsible for creating
# distributable files that can be released by a human or a script.
#
# It can either cross-compile for different platforms using xgo, or simply target a specific
# platform. Set GOOS and GOARCH environment variables to disable xgo and instead
# compile locally for a specific platform.
#
# use:
# ./scripts/build.sh 0.4.8
# GOOS=linux GOARCH=amd64 ./scripts/build.sh 0.4.8
set -ex

dir=$(dirname "${BASH_SOURCE[0]}")

version=$1
projectDir="$dir/../.."
basedir="$projectDir/pkg/cli"
outputDir="$projectDir/build/cli"

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

goVersion=go-1.25.x

get_binary_name() {
  platform=$1

  if [ "$platform" == "windows" ]; then
    echo "dnote.exe"
  else
    echo "dnote"
  fi
}

build() {
  platform=$1
  arch=$2
  # native indicates if the compilation is to take place natively on the host platform
  # if not true, use xgo with Docker to cross-compile
  native=$3

  # build binary
  destDir="$outputDir/$platform-$arch"
  ldflags="-X main.apiEndpoint=https://localhost:3001/api -X main.versionTag=$version"
  tags="fts5"

  pushd "$projectDir"

  mkdir -p "$destDir"

  if [ "$native" == true ]; then
    GOOS="$platform" GOARCH="$arch" \
      go build \
        -ldflags "$ldflags" \
        --tags "$tags" \
        -o="$destDir/cli-$platform-$arch" \
        "$basedir"
  else
    flags=()
    if [ "$platform" == "windows" ]; then
      flags+=("-buildmode=exe")
    fi

    xgo \
      -go "$goVersion" \
      -targets="$platform/$arch" \
      -ldflags "$ldflags" \
      -dest="$destDir" \
      -out="cli" \
      "${flags[@]}" \
      -tags "$tags" \
      -pkg pkg/cli \
      .
  fi

  popd

  binaryName=$(get_binary_name "$platform")
  mv "$destDir/cli-"* "$destDir/$binaryName"

  # build tarball
  tarballName="dnote_${version}_${platform}_${arch}.tar.gz"
  tarballPath="$outputDir/$tarballName"

  cp "$projectDir/LICENSE" "$destDir"
  cp "$basedir/README.md" "$destDir"
  tar -C "$destDir" -zcvf "$tarballPath" "."
  rm -rf "$destDir"

  # calculate checksum
  pushd "$outputDir"
  shasum -a 256 "$tarballName" >> "$outputDir/dnote_${version}_checksums.txt"
  popd
}

if [ -z "$GOOS" ] && [ -z "$GOARCH" ]; then
  # install the tool
  go install src.techknowlogick.com/xgo@latest

  # Linux
  build linux amd64
  build linux arm64
  build linux arm

  # macOS
  build darwin amd64
  build darwin arm64

  # Windows
  build windows amd64

  # FreeBSD
  build freebsd amd64
else
  build "$GOOS" "$GOARCH" true
fi
