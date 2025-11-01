/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package crypt

import (
	"crypto/rand"

	"encoding/base64"
	"github.com/pkg/errors"
)

// getRandomBytes generates a cryptographically secure pseudorandom numbers of the
// given size in byte
func getRandomBytes(numBytes int) ([]byte, error) {
	b := make([]byte, numBytes)
	if _, err := rand.Read(b); err != nil {
		return nil, errors.Wrap(err, "reading random bits")
	}

	return b, nil
}

// GetRandomStr generates a cryptographically secure pseudorandom numbers of the
// given size in byte
func GetRandomStr(numBytes int) (string, error) {
	b, err := getRandomBytes(numBytes)
	if err != nil {
		return "", errors.Wrap(err, "generating random bits")
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
