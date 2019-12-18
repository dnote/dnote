package token

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// generateRandom generates random bits of given length
func generateRandom(bits int) (string, error) {
	b := make([]byte, bits)

	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "reading random bytes")
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// Create generates a new token in the database
func Create(db *gorm.DB, userID int, kind string) (database.Token, error) {
	val, err := generateRandom(16)
	if err != nil {
		return database.Token{}, errors.Wrap(err, "generating random bytes")
	}

	token := database.Token{
		UserID: userID,
		Value:  val,
		Type:   kind,
	}
	if err := db.Save(&token).Error; err != nil {
		return database.Token{}, errors.Wrap(err, "creating a token for unsubscribing")
	}

	return token, nil
}
