#!/bin/bash

set -eu

command_exists () {
  command -v "$1" >/dev/null 2>&1;
}

if [ $# -eq 0 ]; then
  echo "no version specified."
  exit 1
fi
if [[ $1 == v* ]]; then 
  echo "do not prefix version with v"
  exit 1
fi

if ! command_exists hub; then
  echo "please install hub"
  exit 1
fi
if ! command_exists shasum; then
  echo "please install shasum"
  exit 1
fi

binary=dnote
version=$1
version_tag="v$version"
goos=("linux" "openbsd" "freebsd" "darwin" "windows")
goarch=("386" "amd64")

rm -rf ./release
mkdir ./release
cp LICENSE ./release/LICENSE
cp README.md ./release/README.md

echo "* release $version"

# 1. build
for os in "${goos[@]}"; do
  for arch in "${goarch[@]}"; do
    filename="${binary}_${version}_${os}_${arch}"
    echo "* building $filename"

    goos="$os" goarch="$arch" go build \
      -o "./release/$filename" \
      -ldflags "-X main.apiEndpoint=https://api.dnote.io -X main.versionTag=$version"

    pushd ./release > /dev/null
    cp "$filename" dnote
    tar -czvf "$filename.tar.gz" dnote LICENSE README.md
    shasum -a 256 "$filename" >> "dnote_${version}_checksums.txt"
    popd > /dev/null
  done
done

# 2. push tag
echo "* tagging and pushing the tag"
git tag -a "$version_tag" -m "Release $version_tag"
git push --tags

# 3. create release
files=(./release/*.tar.gz ./release/*.txt)
file_args=()
for file in "${files[@]}"; do
  file_args+=("--attach=$file")
done

echo "* creating release"
set -x
hub release create \
  "${file_args[@]}" \
  --message="$version_tag"\
  "$version_tag"

# 4. release on brew

homebrew_sha256=$(shasum -a 256 "./release/${binary}_${version}_darwin_amd64.tar.gz" | cut -d ' ' -f 1)
(cd "$GOPATH"/src/github.com/dnote/homebrew-dnote && ./release.sh "$version" "$homebrew_sha256")
