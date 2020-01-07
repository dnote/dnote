/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package app

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/dnote/dnote/pkg/server/testutils"
)

func TestMarkDigestRead(t *testing.T) {
	defer testutils.ClearData()

	user := testutils.SetupUserData()
	digest := database.Digest{UserID: user.ID}
	testutils.MustExec(t, testutils.DB.Save(&digest), "preparing digest")

	a := NewTest(nil)

	// Multiple calls should not create more than 1 receipt
	for i := 0; i < 3; i++ {
		ret, err := a.MarkDigestRead(digest, user)
		if err != nil {
			t.Fatal(err, "failed to perform")
		}

		var receiptCount int
		testutils.MustExec(t, testutils.DB.Model(&database.DigestReceipt{}).Count(&receiptCount), "counting receipts")
		assert.Equalf(t, receiptCount, 1, "receipt count mismatch")

		var receipt database.DigestReceipt
		testutils.MustExec(t, testutils.DB.Where("id = ?", ret.ID).First(&receipt), "getting receipt")
		assert.Equalf(t, receipt.UserID, user.ID, "receipt UserID mismatch")
		assert.Equalf(t, receipt.DigestID, digest.ID, "receipt DigestID mismatch")
	}
}
