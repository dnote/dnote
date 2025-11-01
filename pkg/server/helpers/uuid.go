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

package helpers

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// GenUUID generates a new uuid v4
func GenUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "generating uuid")
	}

	return id.String(), nil
}

// ValidateUUID validates the given uuid
func ValidateUUID(u string) bool {
	_, err := uuid.Parse(u)

	return err == nil
}
