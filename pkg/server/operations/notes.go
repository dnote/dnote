package operations

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/helpers"
	"github.com/dnote/dnote/pkg/server/permissions"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GetNote retrieves a note for the given user
func GetNote(db *gorm.DB, uuid string, user database.User) (database.Note, bool, error) {
	zeroNote := database.Note{}
	if !helpers.ValidateUUID(uuid) {
		return zeroNote, false, nil
	}

	conn := db.Where("notes.uuid = ? AND deleted = ?", uuid, false)
	conn = database.PreloadNote(conn)

	var note database.Note
	conn = conn.Find(&note)

	if conn.RecordNotFound() {
		return zeroNote, false, nil
	} else if err := conn.Error; err != nil {
		return zeroNote, false, errors.Wrap(err, "finding note")
	}

	if ok := permissions.ViewNote(&user, note); !ok {
		return zeroNote, false, nil
	}

	return note, true, nil
}
