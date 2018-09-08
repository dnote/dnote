# Contributing

This is a guide for contributors.

## Set up

First, download the project

```sh
go get github.com/dnote/cli
```

Go to the project root and download dependencies using [dep](https://github.com/golang/dep).

```sh
dep ensure
```

## Test

Run

```sh
./scripts/test.sh
```

## Debug

Run Dnote with `DNOTE_DEBUG=1` to print debugging statements.

## Release

This project uses [goreleaser](https://github.com/goreleaser/goreleaser) to automate the release process.

The following will tag, push the tag, create release on GitHub, build artifacts, upload them, and
push a commit to [Dnote Homebrew tap](https://github.com/dnote/homebrew-dnote).

```sh
VERSION=v0.4.2 make
```

**Note**

- If a release is not stable,
  - disable the homebrew release by commenting out `homebrew` block in `.goreleaser.yml`
  - mark release as pre-release on GitHub release
