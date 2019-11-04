package permissions

import (
	"github.com/dnote/dnote/pkg/server/database"
)

// ViewNote checks if the given user can view the given note
func ViewNote(user *database.User, note database.Note) bool {
	if note.Public {
		return true
	}
	if user == nil {
		return false
	}

	return note.UserID == user.ID
}
