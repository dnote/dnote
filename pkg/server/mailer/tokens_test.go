/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
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
