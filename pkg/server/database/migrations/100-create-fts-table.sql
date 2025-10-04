-- Create FTS5 virtual table for full-text search on notes
CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(
    content=notes,
    body,
    tokenize="porter unicode61 categories 'L* N* Co Ps Pe'"
);

-- Create triggers to keep notes_fts in sync with notes
CREATE TRIGGER IF NOT EXISTS notes_insert AFTER INSERT ON notes BEGIN
    INSERT INTO notes_fts(rowid, body) VALUES (new.rowid, new.body);
END;
CREATE TRIGGER IF NOT EXISTS notes_delete AFTER DELETE ON notes BEGIN
    INSERT INTO notes_fts(notes_fts, rowid, body) VALUES ('delete', old.rowid, old.body);
END;
CREATE TRIGGER IF NOT EXISTS notes_update AFTER UPDATE ON notes BEGIN
    INSERT INTO notes_fts(notes_fts, rowid, body) VALUES ('delete', old.rowid, old.body);
    INSERT INTO notes_fts(rowid, body) VALUES (new.rowid, new.body);
END;