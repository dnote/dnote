# Dnote CLI

![Dnote](assets/main.png)

A command line interface for spontaneously capturing the things you learn while coding

## Installation

On macOS, or Linux, run:

    curl -s https://raw.githubusercontent.com/dnote-io/cli/master/install.sh | sh

In some cases, you might need `sudo`. Feel free to inspect [install.sh](https://github.com/dnote-io/cli/blob/master/install.sh):

    curl -s https://raw.githubusercontent.com/dnote-io/cli/master/install.sh | sudo sh

On Windows, download [binary](https://github.com/dnote-io/cli/releases)

## Usage

Dnote categorizes your **notes** by **books**.

All your books and notes are stored in `$HOME/.dnote` as a YAML file.

In the future, you can sync your note with the Dnote server and set up digest notifications to reinforce your learning.

### Commands

**dnote use [book name]**

Change the book to write your note in.

e.g.

    dnote use linux

**dnote new "[note]"**

Write a new note under the current book.

e.g.

    dnote new "set -e instructs bash to exit immediately if any command has non-zero exit status"

**dnote books**

List all the books that you created

e.g.

    $ dnote books
      javascript
    * linux
      tmux
      css


## Links

* [Website](https://dnote.io)
* [Making Dnote (blog article)](https://sungwoncho.io/making-dnote/)

## License

MIT

-------

> Made by [sung](https://sungwoncho.io)
