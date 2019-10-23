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

package repetition

import (
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func getRuleBookIDs(ruleID int) ([]int, error) {
	db := database.DBConn
	var ret []int
	if err := db.Table("repetition_rule_books").Select("book_id").Where("repetition_rule_id = ?", ruleID).Pluck("book_id", &ret).Error; err != nil {
		return nil, errors.Wrap(err, "querying book_ids")
	}

	return ret, nil
}

func applyBookDomain(noteQuery *gorm.DB, rule database.RepetitionRule) (*gorm.DB, error) {
	ret := noteQuery

	if rule.BookDomain != database.BookDomainAll {
		bookIDs, err := getRuleBookIDs(rule.ID)
		if err != nil {
			return nil, errors.Wrap(err, "getting book_ids")
		}

		ret = ret.Joins("INNER JOIN books ON notes.book_uuid = books.uuid")

		if rule.BookDomain == database.BookDomainExluding {
			ret = ret.Where("books.id NOT IN (?)", bookIDs)
		} else if rule.BookDomain == database.BookDomainIncluding {
			ret = ret.Where("books.id IN (?)", bookIDs)
		}
	}

	return ret, nil
}

// getRandomNotes returns a random set of notes
func getRandomNotes(db *gorm.DB, rule database.RepetitionRule) ([]database.Note, error) {
	conn := db.Table("notes").Where("notes.user_id = ?", rule.UserID)
	conn, err := applyBookDomain(conn, rule)
	if err != nil {
		return nil, errors.Wrap(err, "applying book domain")
	}

	// TODO: Find alternatives because ordering by random() does not scale with the number of rows
	conn = conn.Order("random()").Limit(rule.NoteCount)

	var notes []database.Note
	if err := conn.Preload("Book").Find(&notes).Error; err != nil {
		return notes, errors.Wrap(err, "getting notes")
	}

	return notes, nil
}
