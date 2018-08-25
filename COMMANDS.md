# Commands

- [add](#dnote-add)
- [view](#dnote-view)
- [edit](#dnote-edit)
- [remove](#dnote-remove)
- [upgrade](#dnote-upgrade)
- [login](#dnote-login)
- [sync](#dnote-sync)

## dnote add

_alias: a, n, new_

Add a new note to a book.

### `dnote add [book name]`

Launch a text editor to add a new note to the specified book.

### `dnote add [book name] -c "[content]"`

Write a new note with a content to the specified book.

e.g.

    $ dnote add linux -c "find - recursively walk the directory"

## dnote view

_alias: v_

- List books or notes.
- View a note detail.

### `dnote view`

List all books.

### `dnote view [book name]`

List all notes in the book.

e.g

    $ dnote view
    $ dnote view golang

### `dnote view [book name] [note index]`

See details of a note

e.g

    $ dnote view golang 12

## dnote edit

_alias: e_

Edit a note

### `dnote edit [book name] [note index]`

Launch a text editor to edit a note with the given index.

### `dnote edit [book name] [note index] -c "[note content]"`

Edit a note with the given index in the specified book with a content.

e.g

    $ dnote edit linux 1 "New Content"

## dnote remove

_alias: d_

Remove either a note or a book

### `dnote remove [book name] [index]`

Removes the note with `index` in the specified book.

### `dnote remove -b [book name]`

Removes the book with the `book name`.

e.g

    $ dnote remove JS 1
    $ dnote remove -b JS

## dnote upgrade

Upgrade the Dnote if newer release is available

## dnote sync

_Dnote Cloud only_

Sync notes with Dnote cloud

## dnote login

_Dnote Cloud only_

Start a login prompt

---

# Deprecated Commands

The following commands are deprecated and easier alternatives are provided.

## dnote ls

**"ls" is replaced by "view". It will be removed in v0.5.0**

_alias: l, notes_

List books or notes

### `dnote ls`

List all books.

### `dnote ls [book name]`

List all notes in the book.

e.g

    $ dnote ls
    $ dnote ls golang

## dnote cat

**"cat" is replaced by "view". It will be removed in v0.5.0**

_alias: c_

See details of a note

### `dnote cat [book name] [note index]`

e.g

    $ dnote cat golang 12
