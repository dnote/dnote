package operations

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// CreateDigest creates a new digest
func CreateDigest(db *gorm.DB, rule database.RepetitionRule, notes []database.Note) (database.Digest, error) {
	var maxVersion int
	if err := db.Raw("SELECT COALESCE(max(version), 0) FROM digests WHERE rule_id = ?", rule.ID).Row().Scan(&maxVersion); err != nil {
		return database.Digest{}, errors.Wrap(err, "finding max version")
	}

	digest := database.Digest{
		RuleID:  rule.ID,
		UserID:  rule.UserID,
		Version: maxVersion + 1,
		Notes:   notes,
	}
	if err := db.Save(&digest).Error; err != nil {
		return database.Digest{}, errors.Wrap(err, "saving digest")
	}

	return digest, nil
}
