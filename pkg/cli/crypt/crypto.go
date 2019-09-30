/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package crypt provides cryptographic funcitonalities
package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"

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
func MakeKeys(password, email []byte, iteration int) ([]byte, []byte, error) {
	masterKey := pbkdf2.Key([]byte(password), []byte(email), iteration, 32, sha256.New)

	authKey, err := runHkdf(masterKey, email, []byte("auth"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "deriving auth key")
	}

	return masterKey, authKey, nil
}

// AesGcmEncrypt encrypts the plaintext using AES in a GCM mode. It returns
// a ciphertext prepended by a 12 byte pseudo-random nonce, encoded in base64.
func AesGcmEncrypt(key, plaintext []byte) (string, error) {
	if key == nil {
		return "", errors.New("no key provided")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "initializing aes")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "initializing gcm")
	}

	nonce := make([]byte, aesGcmNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "generating nonce")
	}

	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
	cipherKeyB64 := base64.StdEncoding.EncodeToString(ciphertext)

	return cipherKeyB64, nil
}

// AesGcmDecrypt decrypts the encrypted data using AES in a GCM mode. The data should be
// a base64 encoded string in the format of 12 byte nonce followed by a ciphertext.
func AesGcmDecrypt(key []byte, dataB64 string) ([]byte, error) {
	if key == nil {
		return nil, errors.New("no key provided")
	}

	data, err := base64.StdEncoding.DecodeString(dataB64)
	if err != nil {
		return nil, errors.Wrap(err, "decoding base64 data")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "initializing aes")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "initializing gcm")
	}

	if len(data) < aesGcmNonceSize {
		return nil, errors.Wrap(err, "malformed data")
	}

	nonce, ciphertext := data[:aesGcmNonceSize], data[aesGcmNonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "decrypting")
	}

	return plaintext, nil
}
