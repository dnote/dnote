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

```bash
# Linux, macOS, FreeBSD, Windows
curl -s https://www.getdnote.com/install | sh

# macOS with Homebrew
brew install dnote
```

Or [download binary](https://github.com/dnote/dnote/releases).

## Server (Optional)

Server is a binary with SQLite embedded. No database setup is required.

If using docker, create a compose.yml:

```yaml
services:
  dnote:
    image: dnote/dnote:latest
    container_name: dnote
    ports:
      - 3001:3001
    volumes:
      - ./dnote_data:/data
    restart: unless-stopped
```

Then run:

```bash
docker-compose up -d
```

Or see the [guide](https://www.getdnote.com/docs/server/manual) for binary installation.

## Documentation

See the [Dnote doc](https://www.getdnote.com/docs).
