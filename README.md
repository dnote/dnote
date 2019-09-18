![Dnote](assets/logo.png)
=========================

Dnote is a simple notebook for developers.

[![Build Status](https://travis-ci.org/dnote/dnote.svg?branch=master)](https://travis-ci.org/dnote/dnote)

## What is Dnote?

Dnote is a lightweight notebook for writing technical notes and neatly organizing them into books. The main design goal is to **keep you focused** by providing a way of swiftly capturing new information **without having to switch environment**. To that end, you can use Dnote as a command line interface, browser extension, web client, or an IDE plugin.

It also offers **end-to-end encrypted** backup with AES-256, a seamless **multi device sync**, and **automated spaced repetition** to retain your memory in case you are building a personal knowledge base.

For more details, see the [download page](https://www.getdnote.com/download) and [features](https://www.getdnote.com/pricing).

![A demo of Dnote CLI](assets/cli.gif)

## Quick install

The quickest way to try Dnote is to install the command line interface.

### Install with Homebrew

On macOS, you can install using Homebrew:

```sh
brew tap dnote/dnote
brew install dnote

# to upgrade to the latest version
brew upgrade dnote
```

### Install with script

You can use the installation script to install the latest version:

    curl -s https://raw.githubusercontent.com/dnote/dnote/master/pkg/cli/install.sh | sh

In some cases, you might need an elevated permission:

    curl -s https://raw.githubusercontent.com/dnote/dnote/master/pkg/cli/install.sh | sudo sh

### Install with tarball

You can download the binary for your platform manually from the [releases page](https://github.com/dnote/dnote/releases).

## Personal knowledge base

Dnote is great for building a personal knowledge base because:

* It is fully open source.
* Your data is stored locally first and in a SQLite format which is [suitable for continued accessibility](https://www.sqlite.org/locrsf.html).
* It provides a way of instantly capturing new lessons without distracting you.
* It automates spaced repetition to help you retain your memory.

You can read more in the following user stories:

- [How I Built a Personal Knowledge Base for Myself](https://www.getdnote.com/blog/how-i-built-personal-knowledge-base-for-myself/)
- [I Wrote Down Everything I Learned While Programming for a Month](https://www.getdnote.com/blog/writing-everything-i-learn-coding-for-a-month/)

## See Also

- [Homepage](https://www.getdnote.com)
- [Forum](https://forum.getdnote.com)
