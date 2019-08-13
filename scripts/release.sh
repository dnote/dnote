#!/usr/bin/env bash
#
# release.sh releases the tarballs and checksum in the build directory
# to GitHub and brew. A prerequisite is to build those files using build.sh.
# use: ./scripts/release.sh cli v0.4.8 path/to/assets

set -euxo pipefail

project=$1
version=$2
assetPath=$3

if [ "$project" != "cli" ] && [ "$project" != "server" ]; then
  echo "unrecognized project '$project'"
  exit 1
fi
if [ -z "$version" ]; then
  echo "no version specified."
  exit 1
fi
if [[ $version == v* ]]; then
  echo "do not prefix version with v"
  exit 1
fi

# 1. push tag
version_tag="$project-v$version"

echo "* tagging and pushing the tag"
git tag -a "$version_tag" -m "Release $version_tag"
git push --tags

# 2. release on GitHub
files=("$assetPath"/*)
file_flags=()
for file in "${files[@]}"; do
  file_flags+=("--attach=$file")
done

# mark as prerelease if version is not in a form of major.minor.patch
# e.g. 1.0.1-beta.1
flags=()
if [[ ! "$version" =~ ^[0-9]+.[0-9]+.[0-9]+$ ]]; then
  flags+=("--prerelease")
fi

echo "* creating release"
set -x
hub release create \
  "${file_flags[@]}" \
  "${flags[@]}" \
  --message="$version_tag"\
  "$version_tag"
