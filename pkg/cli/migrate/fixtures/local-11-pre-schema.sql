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
CREATE UNIQUE INDEX idx_books_uuid ON books(uuid);
CREATE TABLE IF NOT EXISTS "notes"
		(
			uuid text NOT NULL,
			book_uuid text NOT NULL,
			body text NOT NULL,
			added_on integer NOT NULL,
			edited_on integer DEFAULT 0,
			public bool DEFAULT false,
			dirty bool DEFAULT false,
			usn int DEFAULT 0 NOT NULL,
			deleted bool DEFAULT false
		);
CREATE VIRTUAL TABLE note_fts USING fts5(content=notes, body, tokenize="porter unicode61 categories 'L* N* Co Ps Pe'")
/* note_fts(body) */;
CREATE TABLE IF NOT EXISTS 'note_fts_data'(id INTEGER PRIMARY KEY, block BLOB);
CREATE TABLE IF NOT EXISTS 'note_fts_idx'(segid, term, pgno, PRIMARY KEY(segid, term)) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS 'note_fts_docsize'(id INTEGER PRIMARY KEY, sz BLOB);
CREATE TABLE IF NOT EXISTS 'note_fts_config'(k PRIMARY KEY, v) WITHOUT ROWID;
CREATE TRIGGER notes_after_insert AFTER INSERT ON notes BEGIN
				INSERT INTO note_fts(rowid, body) VALUES (new.rowid, new.body);
			END;
CREATE TRIGGER notes_after_delete AFTER DELETE ON notes BEGIN
				INSERT INTO note_fts(note_fts, rowid, body) VALUES ('delete', old.rowid, old.body);
			END;
CREATE TRIGGER notes_after_update AFTER UPDATE ON notes BEGIN
				INSERT INTO note_fts(note_fts, rowid, body) VALUES ('delete', old.rowid, old.body);
				INSERT INTO note_fts(rowid, body) VALUES (new.rowid, new.body);
			END;
CREATE TABLE actions
		(
			uuid text PRIMARY KEY,
			schema integer NOT NULL,
			type text NOT NULL,
			data text NOT NULL,
			timestamp integer NOT NULL
		);
CREATE UNIQUE INDEX idx_notes_uuid ON notes(uuid);
CREATE INDEX idx_notes_book_uuid ON notes(book_uuid);
