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

* Build for all target platforms, tag, push tags
* Release on GitHub and [Dnote Homebrew tap](https://github.com/dnote/homebrew-dnote).

```sh
VERSION=0.4.8 make release
```

* Build, without releasing, for all target platforms

```sh
VERSION=0.4.8 make build
```

**Note**

- If a release is not stable,
  - disable the homebrew release by commenting out `homebrew` block in `.goreleaser.yml`
  - mark release as pre-release on GitHub release
