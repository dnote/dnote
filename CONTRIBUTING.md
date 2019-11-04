# Contributing to Dnote

Dnote is an open source project.

* [Setting up](#setting-up)
* [Server](#server)
* [Command Linux Interface](#command-line-interface)

## Setting up

Dnote uses [Vagrant](https://github.com/hashicorp/vagrant) to provision a consistent development environment.

*Prerequisites*

* Vagrant ([Download](https://www.vagrantup.com/downloads.html))
* VirtualBox ([Download](https://www.virtualbox.org/))

Following steps will set up your development environment and install dependencies in a virtual machine.

1. Run `vagrant up` to start a virtual machine and bootstrap the development environment.
2. Run `vagrant rsync-auto` to sync the files with the virtual machine.

*Workflow*

* You can make changes to the source code from the host machine.
* Any commands need to be run inside the virtual machine. You can connect to it by running `vagrant ssh`.

## Server

The server consists of the frontend web application and a web server.

### Development

* Run `make dev-server` to start a local server.
* You can access the server on `localhost:3000` on your machine.

### Test

```bash
# Run tests for the frontend web application
make test-web

# Run tests for API
make test-api
```


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

