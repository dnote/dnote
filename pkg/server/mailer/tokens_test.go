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

package mailer

import (
	"testing"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestGetToken(t *testing.T) {
	db := testutils.InitMemoryDB(t)

	userID := 1
	tokenType := "email_verification"

	t.Run("creates new token", func(t *testing.T) {
		token, err := GetToken(db, userID, tokenType)
		if err != nil {
			t.Fatalf("GetToken failed: %v", err)
		}

		if token.UserID != userID {
			t.Errorf("expected UserID %d, got %d", userID, token.UserID)
		}
		if token.Type != tokenType {
			t.Errorf("expected Type %s, got %s", tokenType, token.Type)
		}
		if token.Value == "" {
			t.Error("expected non-empty token Value")
		}
		if token.UsedAt != nil {
			t.Error("expected UsedAt to be nil for new token")
		}
	})

	t.Run("reuses unused token", func(t *testing.T) {
		// Get token again - should return the same one
		token2, err := GetToken(db, userID, tokenType)
		if err != nil {
			t.Fatalf("second GetToken failed: %v", err)
		}

		// Get first token to compare
		var token1 database.Token
		if err := db.Where("user_id = ? AND type = ?", userID, tokenType).First(&token1).Error; err != nil {
			t.Fatalf("failed to get first token: %v", err)
		}

		if token1.ID != token2.ID {
			t.Errorf("expected same token ID %d, got %d", token1.ID, token2.ID)
		}
		if token1.Value != token2.Value {
			t.Errorf("expected same token Value %s, got %s", token1.Value, token2.Value)
		}

		// Verify only one token exists in database
		var count int64
		if err := db.Model(&database.Token{}).Where("user_id = ? AND type = ?", userID, tokenType).Count(&count).Error; err != nil {
			t.Fatalf("failed to count tokens: %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 token in database, got %d", count)
		}
	})
}
