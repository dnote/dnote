#!/bin/bash
#
# release.sh releases the tarballs and checksum in the build directory
# to GitHub and brew. It is important to build those files using build.sh
# use: ./scripts/release.sh v0.4.8

set -eu

homebrewRepoDir="$GOPATH"/src/github.com/dnote/homebrew-dnote

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

if [ ! -d "$homebrewRepoDir" ]; then
  echo "homebrew-dnote not found locally. did you clone it?"
  exit 1
fi

# 1. push tag
version=$1
version_tag="cli-v$version"

echo "* tagging and pushing the tag"
git tag -a "$version_tag" -m "Release $version_tag"
git push --tags

# 2. release on GitHub
files=(./build/*.tar.gz ./build/*.txt)
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

# 3. Release on Homebrew
homebrew_sha256=$(shasum -a 256 "./build/dnote_${version}_darwin_amd64.tar.gz" | cut -d ' ' -f 1)
(cd "$homebrewRepoDir" && ./release.sh "$version" "$homebrew_sha256")
