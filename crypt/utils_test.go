package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/dnote/cli/testutils"
	"github.com/pkg/errors"
)

func TestAesGcmEncrypt(t *testing.T) {
	testCases := []struct {
		key       []byte
		plaintext []byte
	}{
		{
			key:       []byte("AES256Key-32Characters1234567890"),
			plaintext: []byte("foo bar baz quz"),
		},
		{
			key:       []byte("AES256Key-32Charactersabcdefghij"),
			plaintext: []byte("1234 foo 5678 bar 7890 baz"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("key %s plaintext %s", tc.key, tc.plaintext), func(t *testing.T) {
			// encrypt
			data, err := AesGcmEncrypt(tc.key, tc.plaintext)
			if err != nil {
				t.Fatal(errors.Wrap(err, "performing encryption"))
			}

			nonce, ciphertext := data[:12], data[12:]

			fmt.Println(string(data))

			block, err := aes.NewCipher([]byte(tc.key))
			if err != nil {
				t.Fatal(errors.Wrap(err, "initializing aes"))
			}

			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				t.Fatal(errors.Wrap(err, "initializing gcm"))
			}

			plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
			if err != nil {
				t.Fatal(errors.Wrap(err, "decode"))
			}

			testutils.AssertDeepEqual(t, plaintext, tc.plaintext, "plaintext mismatch")
		})
	}
}

func TestAesGcmDecrypt(t *testing.T) {
	testCases := []struct {
		key               []byte
		ciphertextB64     string
		expectedPlaintext []byte
	}{
		{
			key:               []byte("AES256Key-32Characters1234567890"),
			ciphertextB64:     "M2ov9hWMQ52v1S/zigwX3bJt4cVCV02uiRm/grKqN/rZxNkJrD7vK4Ii0g==",
			expectedPlaintext: []byte("foo bar baz quz"),
		},
		{
			key:               []byte("AES256Key-32Characters1234567890"),
			ciphertextB64:     "M4csFKUIUbD1FBEzLgHjscoKgN0lhMGJ0n2nKWiCkE/qSKlRP7kS",
			expectedPlaintext: []byte("foo\n1\nbar\n2"),
		},
		{
			key:               []byte("AES256Key-32Characters1234567890"),
			ciphertextB64:     "pe/fnw73MR1clmVIlRSJ5gDwBdnPly/DF7DsR5dJVz4dHZlv0b10WzvJEGOCHZEr+Q==",
			expectedPlaintext: []byte("föo\nbār\nbåz & qūz"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("key %s ciphertext %s", tc.key, tc.ciphertextB64), func(t *testing.T) {
			ciphertext, err := base64.StdEncoding.DecodeString(tc.ciphertextB64)
			if err != nil {
				t.Fatal(errors.Wrap(err, "decoding ciphertext base64"))
			}

			plaintext, err := AesGcmDecrypt(tc.key, ciphertext)
			if err != nil {
				t.Fatal(errors.Wrap(err, "performing decryption"))
			}

			testutils.AssertDeepEqual(t, plaintext, tc.expectedPlaintext, "plaintext mismatch")
		})
	}
}
