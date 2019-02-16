// Package crypt provides cryptographic funcitonalities
package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

var aesGcmNonceSize = 12

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

	encKeyB64 := base64.StdEncoding.EncodeToString(encKey)
	authKeyB64 := base64.StdEncoding.EncodeToString(authKey)

	return encKeyB64, authKeyB64, nil
}

// AesGcmEncrypt encrypts the plaintext using AES in a GCM mode. It returns
// a ciphertext prepended by a 12 byte pseudo-random nonce.
func AesGcmEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, errors.Wrap(err, "initializing aes")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, errors.Wrap(err, "initializing gcm")
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	return aesgcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

// AesGcmDecrypt decrypts the encrypted data using AES in a GCM mode.
func AesGcmDecrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, errors.Wrap(err, "initializing aes")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, errors.Wrap(err, "initializing gcm")
	}

	nonce, ciphertext := data[:aesGcmNonceSize], data[aesGcmNonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "decrypting")
	}

	return plaintext, nil
}
