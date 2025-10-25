package database

import (
	"testing"
)

func TestSchemaEmbed(t *testing.T) {
	db := InitTestMemoryDB(t)
	defer db.Close()

	// Try to insert a book to verify schema is loaded
	_, err := db.Exec("INSERT INTO books (uuid, label) VALUES (?, ?)", "test-uuid", "test-label")
	if err != nil {
		t.Fatalf("Failed to insert into books: %v", err)
	}

	// Verify it was inserted
	var label string
	err = db.QueryRow("SELECT label FROM books WHERE uuid = ?", "test-uuid").Scan(&label)
	if err != nil {
		t.Fatalf("Failed to query books: %v", err)
	}
	if label != "test-label" {
		t.Fatalf("Expected label 'test-label', got '%s'", label)
	}
	t.Log("Schema embed test passed!")
}
