#!/bin/bash
#
# build.sh compiles dnote binary for target platforms. It is resonsible for creating
# distributable files that can be released by a human or a script.
# use: ./scripts/build.sh 0.4.8

set -eux

version=$1
projectDir="$GOPATH/src/github.com/dnote/dnote"
basedir="$GOPATH/src/github.com/dnote/dnote/pkg/cli"
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

goVersion=1.12.x

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

  # build binary
  destDir="$outputDir/$platform-$arch"

  mkdir -p "$destDir"
  xgo \
    -go "$goVersion" \
    -ldflags "-X main.apiEndpoint=https://api.dnote.io -X main.versionTag=$version" \
    --targets="$platform/$arch" \
    --tags "fts5" \
    --dest="$destDir" \
    "$basedir"

  binaryName=$(get_binary_name "$platform")
  mv "$destDir/cli-${platform}-"* "$destDir/$binaryName"

  # build tarball
  tarballName="dnote_${version}_${platform}_${arch}.tar.gz"
  tarballPath="$outputDir/$tarballName"

  cp "$projectDir/licenses/GPLv3.txt" "$destDir"
  cp "$basedir/README.md" "$destDir"
  tar -C "$destDir" -zcvf "$tarballPath" "."
  rm -rf "$destDir"

  # calculate checksum
  pushd "$outputDir"
  shasum -a 256 "$tarballName" >> "$outputDir/dnote_${version}_checksums.txt"
  popd
}

# fetch tool
go get -u github.com/karalabe/xgo

build linux amd64
build darwin amd64
build windows amd64
