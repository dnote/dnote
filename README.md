# Dnote CLI

![Dnote](assets/main.png)

A command line interface for spontaneously capturing the things you learn while coding

## Installation

On macOS, or Linux, run:

    curl -s https://raw.githubusercontent.com/dnote-io/cli/master/install.sh | sh

In some cases, you might need `sudo`. Feel free to inspect [install.sh](https://github.com/dnote-io/cli/blob/master/install.sh):

    curl -s https://raw.githubusercontent.com/dnote-io/cli/master/install.sh | sudo sh

On Windows, download [binary](https://github.com/dnote-io/cli/releases)

## Overview

Dnote categorizes your **notes** by **books**.

All your books and notes are stored in `$HOME/.dnote` as a JSON file.

You can optionally sync your note with Dnote server. Syncing will allow you to interact with your notes using the web frontend at https://dnote.io, and set up digest notifications.

## Commands

### dnote use [book name]
*alias: u*

Change the book to write your note in.

e.g.

    dnote use linux

### dnote new "[note]"
*alias: n*

Write a new note under the current book.

e.g.

    dnote new "set -e instructs bash to exit immediately if any command has non-zero exit status"

### dnote edit "[note]" "[content]"
*alias: e*

Overwrite a note under the current book

#### Option
*`-b [book]`

Specify the name of the book to edit from

e.g

    $ dnote notes
    * [0] - Content index 0.
    * [1] - Content index 1.
    * [2] - Content index 2.

    $ dnote edit 1 "New content"
    [+] Edited Note : 1 

    $ dnote notes
    * [0] - Content index 0.
    * [1] - New content.
    * [2] - Content index 2.

### dnote delete "[note]"
*alias: d*

Delete a note under the current book

#### Option
*`-n [note]`

Specify the index of the note to be deleted

*`-b [book]`

Specify the name of the book to be deleted

e.g 
    $ dnote delete -n 1
    $ dnote delete -b foo

 
### dnote books
*alias: b*

List all the books that you created

e.g.

    $ dnote books
      javascript
    * linux
      tmux
      css

### dnote notes

List all notes in the current book

#### Options

* `-b [book name]`

Specify the name of the book to read from

e.g.

    $ dnote notes
    On note JS
    * .bind() creates a new function
    * arrow function uses less memory than function with .bind()
    * the time passed to setTimeout is minimum, no guaranteed time

### dnote sync

Sync notes with Dnote server

### dnote login

Start a login procedure which will store the APIKey to communicate with the server


## Links

* [Website](https://dnote.io)
* [Making Dnote (blog article)](https://sungwoncho.io/making-dnote/)

## License

MIT

-------

> Made by [sung](https://sungwoncho.io)
