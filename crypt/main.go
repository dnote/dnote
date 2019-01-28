// Package crypt provides cryptographic funcitonalities
package crypt

import (
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

func runHkdf(secret, salt, info []byte) ([]byte, error) {
	r := hkdf.New(sha256.New, secret, salt, info)

	ret := make([]byte, 32)
	_, err := io.ReadFull(r, ret)
	if err != nil {
		return []byte{}, errors.Wrap(err, "reading key bytes")
	}

	return ret, nil
}

// MakeKeys derives, from the given credential, a key set comprising of an encryption key
// and an authentication key
func MakeKeys(password, email []byte, iteration int) (string, string, error) {
	masterKey := pbkdf2.Key([]byte(password), []byte(email), iteration, 32, sha256.New)
	log.Debug("email: %s, password: %s", email, password)

	encKey, err := runHkdf(masterKey, email, []byte("enc"))
	if err != nil {
		return "", "", errors.Wrap(err, "deriving enc key")
	}

	authKey, err := runHkdf(masterKey, email, []byte("auth"))
	if err != nil {
		return "", "", errors.Wrap(err, "deriving auth key")
	}

	encKeyHex := base64.StdEncoding.EncodeToString(encKey)
	authKeyHex := base64.StdEncoding.EncodeToString(authKey)

	return encKeyHex, authKeyHex, err
}
