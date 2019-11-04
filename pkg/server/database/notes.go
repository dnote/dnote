package database

import (
	"github.com/jinzhu/gorm"
)

// PreloadNote preloads the associations for a notes for the given query
func PreloadNote(conn *gorm.DB) *gorm.DB {
	return conn.Preload("Book").Preload("User")
}
