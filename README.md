![Dnote](assets/logo.png)
=========================

![Build Status](https://github.com/dnote/dnote/actions/workflows/ci.yml/badge.svg)

Dnote is a simple command line notebook. Single binary, no dependencies. Since 2017.

Your notes are stored in **one SQLite file** - portable, searchable, and completely under your control. Optional sync between devices via a self-hosted server with REST API access.

```sh
# Add a note (or omit -c to launch your editor)
dnote add linux -c "Check disk usage with df -h"

# View notes in a book
dnote view linux

# Full-text search
dnote find "disk usage"

# Sync notes
dnote sync
```

## Installation

On Unix-like systems (Linux, FreeBSD, macOS), you can use the installation script:

    curl -s https://www.getdnote.com/install | sh

Or on macOS with Homebrew:

```sh
brew tap dnote/dnote
brew install dnote
```

You can also download the binary for your platform from the [releases page](https://github.com/dnote/dnote/releases).

## Server

Self-host your own Dnote server - just run a binary, no database required. [Download](https://github.com/dnote/dnote/blob/master/SELF_HOSTING.md) or run [with Docker](https://github.com/dnote/dnote/blob/master/host/docker/README.md).

## Documentation

Please see [Dnote wiki](https://github.com/dnote/dnote/wiki) for the documentation.
