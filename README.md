![Dnote](assets/logo.png)
=========================

![Build Status](https://github.com/dnote/dnote/actions/workflows/ci.yml/badge.svg)

Dnote is a simple command line notebook for programmers.

It **keeps you focused** by providing a way of effortlessly capturing and retrieving information **without leaving your terminal**. It also offers a seamless **multi-device sync**.

![A demo of Dnote command line interface](assets/cli.gif "Dnote command line interface")

## Installation

On macOS, you can install using Homebrew:

```sh
brew tap dnote/dnote
brew install dnote
```

On Linux or macOS, you can use the installation script:

    curl -s https://www.getdnote.com/install | sh

Otherwise, you can download the binary for your platform manually from the [releases page](https://github.com/dnote/dnote/releases).
## Simple Command

 dnote view- This command showed all the folders that are currently in the working directory and lists them all out in Alphabetical order.

 dnote view <"folder-name"> - This command allows the program to open the folder access the contents inside the folder.

 dnote view <"number"> - There will be numbers denoting each seperate file that is in the folder and all that the user needs to do is type out the number next to view in order to view the file.

 This application can also zip files into .tar.gz format and uncmpress .tar.gz files as well.

 f- specifies the file type

 z - uses gzip to zip files.

 x - extract files

 v - verbose

 In order to specify the output of a compressed file into a specific location it should be formatted in this way: tar -xvzf <"file-name">.tar.fz -C output_directory

 In order to create a new file in a specified directory the command should be:
 dnote add <"folder-name">

In order to exit text editor ":wp"

dnote edit <"number"> - Allows user to edit file containing that specific ID number.

dnote sync - This command syncs the notes with all the platforms that contain the app.

## Example Commands (Shown in Video)

 dnote view

 dnote view bash - (opens the bash folder)

 dnote view 35 - Views the file with the ID of '35'

 dnote edit 36 - Opened the text editor to edit the file with the ID '36'

 dnote add bash - created a new file in the bash folder. 
## Server

The quickest way to experience the Dnote server is to use [Dnote Cloud](https://app.getdnote.com).

Or you can install it on your server by [using Docker](https://github.com/dnote/dnote/blob/master/host/docker/README.md), or [using a binary](https://github.com/dnote/dnote/blob/master/SELF_HOSTING.md).

## Documentation

Please see [Dnote wiki](https://github.com/dnote/dnote/wiki) for the documentation.

## See Also

- [Homepage](https://www.getdnote.com)
- [I Wrote Down Everything I Learned While Programming for a Month](https://www.getdnote.com/blog/writing-everything-i-learn-coding-for-a-month/)
