#!/usr/bin/env bash
set -eux

currentDir=$(dirname "${BASH_SOURCE[0]}")
cliHomebrewDir=${currentDir}/../../homebrew-dnote

if [ ! -d "$cliHomebrewDir" ]; then
  echo "homebrew-dnote not found locally. Cloning."
  git clone git@github.com:dnote/homebrew-dnote.git "$cliHomebrewDir"
fi

version=$1

echo "version: $version"

# Download source tarball and calculate SHA256
source_url="https://github.com/dnote/dnote/archive/refs/tags/cli-v${version}.tar.gz"
echo "Calculating SHA256 for: $source_url"
sha=$(curl -L "$source_url" | shasum -a 256 | cut -d ' ' -f 1)

pushd "$cliHomebrewDir"

echo "pulling latest dnote-homebrew repo"
git checkout master
git pull origin master

cat > ./Formula/dnote.rb << EOF
class Dnote < Formula
  desc "Simple command line notebook for programmers"
  homepage "https://www.getdnote.com"
  url "https://github.com/dnote/dnote/archive/refs/tags/cli-v${version}.tar.gz"
  sha256 "${sha}"
  license "GPL-3.0"
  head "https://github.com/dnote/dnote.git", branch: "master"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X main.apiEndpoint=https://api.getdnote.com -X main.versionTag=#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags), "-tags", "fts5", "./pkg/cli"
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
