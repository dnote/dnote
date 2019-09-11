# Contributing to Dnote

This repository contains the server side and the client side code for Dnote.

* [Setting up](#setting-up)
* [Command Linux Interface](#command-line-interface)
* [Server](#server)

## Setting up

1. Install the following prerequisites if necessary:

* [Go programming language](https://golang.org/dl/) 1.12+
* [Node.js](https://nodejs.org/) 10.16+
* Postgres 10.9+

2. Get the Dnote code:

```sh
go get github.com/dnote/dnote
```

3. Run `make` to install dependencies

## Command Line Interface

### Build

You can build either a development version or a production version:

```
# Build a development version for your platform and place it in your `PATH`.
make debug=true build-cli

# Build a production version for all platforms
make version=v0.1.0 build-cli

# Build a production version for a specific platform
# Note: You cannot cross-compile using this method because Dnote uses CGO
# and requires the OS specific headers.
GOOS=[insert OS] GOARCH=[insert arch] make version=v0.1.0 build-cli
```

### Test

* Run all tests for the command line interface:

```
make test-cli
```

### Debug

Run Dnote with `DNOTE_DEBUG=1` to print debugging statements. For instance:

```
DNOTE_DEBUG=1 dnote sync
```

### Release

* Run `make version=v0.1.0 release-cli` to achieve the following:
  * Build for all target platforms, create a git tag, push all tags to the repository
  * Create a release on GitHub and [Dnote Homebrew tap](https://github.com/dnote/homebrew-dnote).

**Note**

- If a release is not stable,
  - disable the homebrew release by commenting out relevant code in the release script.
  - mark release as pre-release on GitHub release

## Server

The server consists of the frontend web application and a web server.

### Development

* Create a postgres database by running `createdb -O postgres dnote`
* If the role does not exist, you can create it by running `sudo -u postgres createuser postgres`

* Run `make dev-server` to start a local server

### Test

```bash
# Run tests for the frontend web application
make test-web

# Run tests for API
make test-api
```

