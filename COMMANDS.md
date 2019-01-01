# Commands

- [add](#dnote-add)
- [view](#dnote-view)
- [edit](#dnote-edit)
- [remove](#dnote-remove)
- [find](#dnote-find)
- [login](#dnote-login)
- [sync](#dnote-sync)

## dnote add

_alias: a, n, new_

Add a new note to a book.

```bash
# Launch a text editor to add a new note to the specified book.
dnote add linux

# Write a new note with a content to the specified book.
dnote add linux -c "find - recursively walk the directory"
```

## dnote view

_alias: v_

- List books or notes.
- View a note detail.

```bash
# List all books.
dnote view

# List all notes in a book.
dnote view golang

# See details of a note
dnote view golang 12
```

## dnote edit

_alias: e_

Edit a note

```bash
# Launch a text editor to edit a note with the given index.
dnote edit linux 1

# Edit a note with the given index in the specified book with a content.
dnote edit linux 1 -c "New Content"
```

## dnote remove

_alias: d_

Remove either a note or a book

```bash
# Remove the note with `index` in the specified book.
dnote remove JS 1

# Remove the book with the `book name`.
dnote remove -b JS
```

## dnote find

_alias: f

Find notes by keywords

```bash
# find notes by a keyword
dnote find rpoplpush

# find notes by multiple keywords
dnote find "building a heap"

# find notes within a book
dnote find "merge sort" -b algorithm
```

## dnote sync

_Dnote Cloud only_

_alias: s_

Sync notes with Dnote cloud

## dnote login

_Dnote Cloud only_

Start a login prompt
