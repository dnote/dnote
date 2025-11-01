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

package app

import (
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// incrementUserUSN increment the given user's max_usn by 1
// and returns the new, incremented max_usn
func incrementUserUSN(tx *gorm.DB, userID int) (int, error) {
	// First, get the current max_usn to detect transition from empty server
	var user database.User
	if err := tx.Select("max_usn, full_sync_before").Where("id = ?", userID).First(&user).Error; err != nil {
		return 0, errors.Wrap(err, "getting current user state")
	}

	// If transitioning from empty server (MaxUSN=0) to non-empty (MaxUSN=1),
	// set full_sync_before to current timestamp to force all other clients to full sync
	if user.MaxUSN == 0 && user.FullSyncBefore == 0 {
		currentTime := time.Now().Unix()
		if err := tx.Table("users").Where("id = ?", userID).Update("full_sync_before", currentTime).Error; err != nil {
			return 0, errors.Wrap(err, "setting full_sync_before on empty server transition")
		}
	}

	if err := tx.Table("users").Where("id = ?", userID).Update("max_usn", gorm.Expr("max_usn + 1")).Error; err != nil {
		return 0, errors.Wrap(err, "incrementing user max_usn")
	}

	if err := tx.Select("max_usn").Where("id = ?", userID).First(&user).Error; err != nil {
		return 0, errors.Wrap(err, "getting the updated user max_usn")
	}

	return user.MaxUSN, nil
}
