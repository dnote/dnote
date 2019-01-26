# Dnote CLI

A command line interface for spontaneously capturing the things you learn while coding.

![Dnote](assets/dnote.gif)

## Install

On macOS, you can install using Homebrew:

```sh
brew tap dnote/dnote
brew install dnote

# to upgrade to the latest version
brew upgrade dnote
```

On Linux or macOS, you can use the installation script:

    curl -s https://raw.githubusercontent.com/dnote/cli/master/install.sh | sh

In some cases, you might need an elevated permission:

    curl -s https://raw.githubusercontent.com/dnote/cli/master/install.sh | sudo sh

Otherwise, you can download the binary for your platform manually from the [releases page](https://github.com/dnote/cli/releases).

## Overview

Write technical notes without getting distracted from programming. The reasons are:

- We forget exponentially unless we write down what we learn and come back.
- Ideas cannot be grokked unless we can put them down in clear words.

## Examples

- Add a note to a book named `linux`

```
dnote add linux -c "find - recursively walk the directory"
```

- See the notes in a book

```
dnote view linux

• on book linux
(0) find - recursively walk the directory
```

## Commands

Please refer to [commands](/COMMANDS.md).

## Links

* [Website](https://dnote.io)
* [Making Dnote (blog article)](https://github.com/dnote-io/cli)

## License

MIT

[![Build Status](https://travis-ci.org/dnote/cli.svg?branch=master)](https://travis-ci.org/dnote/cli)
