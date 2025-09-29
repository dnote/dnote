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

package presenters

import (
	"time"

	"github.com/dnote/dnote/pkg/server/database"
)

// EmailPreference is a presented email digest
type EmailPreference struct {
	InactiveReminder bool      `json:"inactive_reminder"`
	ProductUpdate    bool      `json:"product_update"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PresentEmailPreference presents a digest
func PresentEmailPreference(p database.EmailPreference) EmailPreference {
	ret := EmailPreference{
		InactiveReminder: p.InactiveReminder,
		ProductUpdate:    p.ProductUpdate,
		CreatedAt:        FormatTS(p.CreatedAt),
		UpdatedAt:        FormatTS(p.UpdatedAt),
	}

	return ret
}
