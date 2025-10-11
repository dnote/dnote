# Contributing to Dnote

Dnote is an open source project.

* [Setting up](#setting-up)
* [Server](#server)
* [Command Line Interface](#command-line-interface)

## Setting up

The CLI and server are single single binary files with SQLite embedded - no databases to install, no containers to run, no VMs required.

**Prerequisites**

* Go 1.25+ ([Download](https://go.dev/dl/))
* Node.js 18+ ([Download](https://nodejs.org/) - only needed for building frontend assets)

**Quick Start**

1. Clone the repository
2. Install dependencies:
   ```bash
   make install
   ```
3. Start developing! Run tests:
   ```bash
   make test
   ```
   Or start the dev server:
   ```bash
   make dev-server
   ```

That's it. You're ready to contribute.

## Server

```bash
# Start dev server (runs on localhost:3001)
make dev-server

# Run tests
make test-api

# Run tests in watch mode
WATCH=true make test-api
```

## Command Line Interface

```bash
# Run tests
make test-cli

# Build dev version (places in your PATH)
make debug=true build-cli

# Build production version for all platforms
make version=v0.1.0 build-cli

# Build for a specific platform
# Note: You cannot cross-compile using this method because Dnote uses CGO
# and requires the OS specific headers.
GOOS=[insert OS] GOARCH=[insert arch] make version=v0.1.0 build-cli

# Debug mode
DNOTE_DEBUG=1 dnote sync
```
