#!/usr/bin/env bash
set -eux

currentDir=$(dirname "${BASH_SOURCE[0]}")
cliHomebrewDir=${currentDir}/../../homebrew-dnote

if [ ! -d "$cliHomebrewDir" ]; then
  echo "homebrew-dnote not found locally. Cloning."
  git clone git@github.com:dnote/homebrew-dnote.git "$cliHomebrewDir"
fi

version=$1
tarball=$2

echo "version: $version"
echo "tarball: $tarball"

sha=$(shasum -a 256 "$tarball" | cut -d ' ' -f 1)

pushd "$cliHomebrewDir"

echo "pulling latest dnote-homebrew repo"
git checkout master
git pull origin master

cat > ./Formula/dnote.rb << EOF
class Dnote < Formula
  desc "A simple command line notebook for programmers"
  homepage "https://www.getdnote.com"
  url "https://github.com/dnote/dnote/releases/download/cli-v${version}/dnote_${version}_darwin_amd64.tar.gz"
  version "${version}"
  sha256 "${sha}"

  def install
    bin.install "dnote"
  end

  test do
    system "#{bin}/dnote", "version"
  end
end
EOF

git add .
git commit --author="Bot <bot@getdnote.com>" -m "Release ${version}"
git push origin master

popd
