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

package permissions

import (
	"github.com/dnote/dnote/pkg/server/database"
)

// ViewNote checks if the given user can view the given note
func ViewNote(user *database.User, note database.Note) bool {
	if user == nil {
		return false
	}
	if note.UserID == 0 {
		return false
	}

	return note.UserID == user.ID
}
