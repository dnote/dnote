# Contributing to Dnote

This repository contains the server side and the client side code for Dnote.

## Set up

1. Download and setup the [Go programming language](https://golang.org/dl/).
2. Download the project

```sh
go get github.com/dnote/dnote
```

## CLI

### Set up

Download dependencies using [dep](https://github.com/golang/dep).

```sh
dep ensure
```

### Test

Run:

```sh
./cli/scripts/test.sh
```

### Debug

Run Dnote with `DNOTE_DEBUG=1` to print debugging statements.

### Release

* Build for all target platforms, tag, push tags
* Release on GitHub and [Dnote Homebrew tap](https://github.com/dnote/homebrew-dnote).

```sh
VERSION=0.4.8 make release
```

* Build, without releasing, for all target platforms

```sh
VERSION=0.4.8 make
```

**Note**

- If a release is not stable,
  - disable the homebrew release by commenting out relevant code in the release script.
  - mark release as pre-release on GitHub release

## Web

### Set up

Download dependencies using [dep](https://github.com/golang/dep) and npm.

```sh
dep ensure
npm install
```

### Test

Run:

```
npm run test
```

## Server

### Set up

Download dependencies using [dep](https://github.com/golang/dep).

```sh
dep ensure
```

### Test

Run:

```
./server/api/scripts/test-local.sh
```
