# Commands

- [add](#dnote-add)
- [view](#dnote-view)
- [edit](#dnote-edit)
- [remove](#dnote-remove)
- [find](#dnote-find)
- [sync](#dnote-sync)
- [login](#dnote-login)
- [logout](#dnote-logout)

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
dnote view 12
```

## dnote edit

_alias: e_

Edit a note or a book.

```bash
# Launch a text editor to edit a note with the given id.
dnote edit 12

# Edit a note with the given id in the specified book with a content.
dnote edit 12 -c "New Content"

# Launch a text editor to edit a book name.
dnote edit js

# Edit a book name by using a flag.
dnote edit js -n "javascript"
```

## dnote remove

_alias: rm, d_

Remove either a note or a book.

```bash
# Remove a note with an id.
dnote remove 1

# Remove a book with the `book name`.
dnote remove js
```

## dnote find

_alias: f_

Find notes by keywords.

```bash
# find notes by a keyword
dnote find rpoplpush

# find notes by multiple keywords
dnote find "building a heap"

# find notes within a book
dnote find "merge sort" -b algorithm
```

## dnote sync

_Dnote Pro only_

_alias: s_

Sync notes with Dnote server. All your data is encrypted before being sent to the server.

## dnote login

_Dnote Pro only_

Start a login prompt.

## dnote logout

_Dnote Pro only_

Log out of Dnote.
