CREATE TABLE notes
		(
			id integer PRIMARY KEY AUTOINCREMENT,
			uuid text NOT NULL,
			book_uuid text NOT NULL,
			content text NOT NULL,
			added_on integer NOT NULL,
			edited_on integer DEFAULT 0,
			public bool DEFAULT false
		, dirty bool DEFAULT false, usn int DEFAULT 0 NOT NULL, deleted bool DEFAULT false);
CREATE TABLE books
		(
			uuid text PRIMARY KEY,
			label text NOT NULL
		, dirty bool DEFAULT false, usn int DEFAULT 0 NOT NULL, deleted bool DEFAULT false);
CREATE TABLE system
		(
			key string NOT NULL,
			value text NOT NULL
		);
CREATE UNIQUE INDEX idx_books_label ON books(label);
CREATE UNIQUE INDEX idx_notes_uuid ON notes(uuid);
CREATE UNIQUE INDEX idx_books_uuid ON books(uuid);
CREATE UNIQUE INDEX idx_notes_id ON notes(id);
CREATE INDEX idx_notes_book_uuid ON notes(book_uuid);
