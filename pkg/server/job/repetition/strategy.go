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
	"sort"
	"time"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func getRuleBookIDs(db *gorm.DB, ruleID int) ([]int, error) {
	var ret []int
	if err := db.Table("repetition_rule_books").Select("book_id").Where("repetition_rule_id = ?", ruleID).Pluck("book_id", &ret).Error; err != nil {
		return nil, errors.Wrap(err, "querying book_ids")
	}

	return ret, nil
}

func applyBookDomain(db *gorm.DB, noteQuery *gorm.DB, rule database.RepetitionRule) (*gorm.DB, error) {
	ret := noteQuery

	if rule.BookDomain != database.BookDomainAll {
		bookIDs, err := getRuleBookIDs(db, rule.ID)
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

func getNotes(db, conn *gorm.DB, rule database.RepetitionRule, dst *[]database.Note) error {
	c, err := applyBookDomain(db, conn, rule)
	if err != nil {
		return errors.Wrap(err, "building query for book threahold 1")
	}

	// TODO: ordering by random() does not scale if table grows large
	if err := c.Where("notes.user_id = ?", rule.UserID).Order("random()").Limit(rule.NoteCount).Preload("Book").Find(&dst).Error; err != nil {
		return errors.Wrap(err, "getting notes")
	}

	return nil
}

// getBalancedNotes returns a set of notes with a 'balanced' ratio of added_on dates
func getBalancedNotes(db *gorm.DB, rule database.RepetitionRule) ([]database.Note, error) {
	now := time.Now()
	t1 := now.AddDate(0, 0, -3).UnixNano()
	t2 := now.AddDate(0, 0, -7).UnixNano()

	// Get notes into three buckets with different threshold values
	var stage1, stage2, stage3 []database.Note
	if err := getNotes(db, db.Where("notes.added_on > ?", t1), rule, &stage1); err != nil {
		return nil, errors.Wrap(err, "Failed to get notes with threshold 1")
	}
	if err := getNotes(db, db.Where("notes.added_on > ? AND notes.added_on < ?", t2, t1), rule, &stage2); err != nil {
		return nil, errors.Wrap(err, "Failed to get notes with threshold 2")
	}
	if err := getNotes(db, db.Where("notes.added_on < ?", t2), rule, &stage3); err != nil {
		return nil, errors.Wrap(err, "Failed to get notes with threshold 3")
	}

	notes := []database.Note{}

	// pick one from each bucket at a time until the result is filled
	i1 := 0
	i2 := 0
	i3 := 0
	k := 0
	for {
		if i1+i2+i3 >= rule.NoteCount {
			break
		}

		// if there are not enough notes to fill the result, break
		if len(stage1) == i1 && len(stage2) == i2 && len(stage3) == i3 {
			break
		}

		if k%3 == 0 {
			if len(stage1) > i1 {
				i1++
			}
		} else if k%3 == 1 {
			if len(stage2) > i2 {
				i2++
			}
		} else if k%3 == 2 {
			if len(stage3) > i3 {
				i3++
			}
		}

		k++
	}

	notes = append(notes, stage1[:i1]...)
	notes = append(notes, stage2[:i2]...)
	notes = append(notes, stage3[:i3]...)

	sort.SliceStable(notes, func(i, j int) bool {
		n1 := notes[i]
		n2 := notes[j]

		return n1.AddedOn > n2.AddedOn
	})

	return notes, nil
}
