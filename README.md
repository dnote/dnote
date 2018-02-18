# Dnote CLI

A command line interface for spontaneously capturing the things you learn while coding.

![Dnote](assets/dnote.gif)

## Install

On Linux or macOS, run:

    curl -s https://raw.githubusercontent.com/dnote-io/cli/master/install.sh | sh

In some cases, you might need an elevated permission:

    curl -s https://raw.githubusercontent.com/dnote-io/cli/master/install.sh | sudo sh

On Windows, download [binary](https://github.com/dnote-io/cli/releases).

## Overview

We learn new things everyday while programming. And we forget most of them exponentially because don't write them down.

Dnote is designed for writing technical notes without switching context from programming. You can capture your learning as they happen with no overhead.

## Examples

- If you want to add a note to a book named linux

```
$ dnote add linux -c "find - recursively walk the directory"
```

- If you want to see the note list of a book named linux

```
$ dnote ls linux
• on book linux
(0) find - recursively walk the directory
```

## Commands

Please refer to [commands](/COMMANDS.md).

## Links

* [Dnote](https://dnote.io)
* [Dnote Cloud](https://dnote.io/cloud)
* [Making Dnote](https://sung.io/making-dnote/)
* [Pitching Dnote](https://sung.io/pitching-dnote/)

## License

MIT

[![Build Status](https://travis-ci.org/dnote-io/cli.svg?branch=master)](https://travis-ci.org/dnote-io/cli)
