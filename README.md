# Dnote CLI

A command line interface for a simple, encrypted notebook that respects your privacy.


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

Write down notes in total privacy, without leaving the command line.

Try keeping technical notes using Dnote. The reasons are:

- You are not distracted from programming while using Dnote.
- We forget exponentially unless we write down what we learn and come back.
- Ideas cannot be grokked unless we can put them down in clear words.

Sync your notes to Dnote server and receive automated email digests for spaced repetition.

## Security

All your data is encrypted using AES256 when you sync with the server. Dnote has zero knowledge about the contents and cannot decrypt your data.

## Commands

Please refer to [commands](/COMMANDS.md).

## Links

- [Dnote](https://dnote.io)
- [Dnote Pro](https://dnote.io/pricing)
- [Browser extension](https://github.com/dnote/browser-extension)

[![Build Status](https://semaphoreci.com/api/v1/dnote/cli/branches/master/badge.svg)](https://semaphoreci.com/dnote/cli)
