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
	"fmt"

	"github.com/dnote/dnote/pkg/server/database"
	"github.com/pkg/errors"
)

func (a *App) getExistingDigestReceipt(userID, digestID int) (*database.DigestReceipt, error) {
	var ret database.DigestReceipt
	conn := a.DB.Where("user_id = ? AND digest_id = ?", userID, digestID).First(&ret)

	if conn.RecordNotFound() {
		return nil, nil
	}
	if err := conn.Error; err != nil {
		return nil, errors.Wrap(err, "querying existing digest receipt")
	}

	return &ret, nil
}

// GetUserDigestByUUID retrives a digest by the uuid for the given user
func (a *App) GetUserDigestByUUID(userID int, uuid string) (*database.Digest, error) {
	var ret database.Digest
	conn := a.DB.Where("user_id = ? AND uuid = ?", userID, uuid).First(&ret)

	if conn.RecordNotFound() {
		return nil, nil
	}
	if err := conn.Error; err != nil {
		return nil, errors.Wrap(err, "finding digest")
	}

	return &ret, nil
}

// MarkDigestRead creates a new digest receipt. If one already exists for
// the given digest and the user, it is a noop.
func (a *App) MarkDigestRead(digest database.Digest, user database.User) (database.DigestReceipt, error) {
	db := a.DB

	existing, err := a.getExistingDigestReceipt(user.ID, digest.ID)
	if err != nil {
		return database.DigestReceipt{}, errors.Wrap(err, "checking existing digest receipt")
	}
	if existing != nil {
		return *existing, nil
	}

	dat := database.DigestReceipt{
		UserID:   user.ID,
		DigestID: digest.ID,
	}
	if err := db.Create(&dat).Error; err != nil {
		return database.DigestReceipt{}, errors.Wrap(err, "creating digest receipt")
	}

	return dat, nil
}

// GetDigestsParam is the params for getting a list of digests
type GetDigestsParam struct {
	UserID  int
	Status  string
	Offset  int
	PerPage int
	Order   string
}

func (p GetDigestsParam) getSubQuery() string {
	orderClause := p.getOrderClause("digests")

	return fmt.Sprintf(`SELECT
	digests.id AS digest_id,
	digests.created_at AS created_at,
	COUNT(digest_receipts.id) AS receipt_count
FROM digests
LEFT JOIN digest_receipts ON digest_receipts.digest_id = digests.id
WHERE digests.user_id = %d
GROUP BY digests.id, digests.created_at
%s`, p.UserID, orderClause)
}

func (p GetDigestsParam) getSubQueryWhere() string {
	var ret string

	if p.Status == "unread" {
		ret = "WHERE t1.receipt_count = 0"
	} else if p.Status == "read" {
		ret = "WHERE t1.receipt_count > 0"
	}

	return ret
}

func (p GetDigestsParam) getOrderClause(table string) string {
	if p.Order == "" {
		return ""
	}

	return fmt.Sprintf(`ORDER BY %s.%s`, table, p.Order)
}

// CountDigests counts digests with the given user using the given criteria
func (a *App) CountDigests(p GetDigestsParam) (int, error) {
	subquery := p.getSubQuery()
	whereClause := p.getSubQueryWhere()
	query := fmt.Sprintf(`SELECT COUNT(*) FROM (%s) AS t1 %s`, subquery, whereClause)

	result := struct {
		Count int
	}{}
	if err := a.DB.Raw(query).Scan(&result).Error; err != nil {
		return 0, errors.Wrap(err, "running count query")
	}

	return result.Count, nil
}

func (a *App) queryDigestIDs(p GetDigestsParam) ([]int, error) {
	subquery := p.getSubQuery()
	whereClause := p.getSubQueryWhere()
	orderClause := p.getOrderClause("t1")
	query := fmt.Sprintf(`SELECT t1.digest_id FROM (%s) AS t1 %s %s OFFSET ? LIMIT ?;`, subquery, whereClause, orderClause)

	ret := []int{}
	rows, err := a.DB.Raw(query, p.Offset, p.PerPage).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "getting rows")
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return []int{}, errors.Wrap(err, "scanning row")
		}

		ret = append(ret, id)
	}

	return ret, nil
}

// GetDigests queries digests for the given user using the given criteria
func (a *App) GetDigests(p GetDigestsParam) ([]database.Digest, error) {
	IDs, err := a.queryDigestIDs(p)
	if err != nil {
		return nil, errors.Wrap(err, "querying digest IDs")
	}

	var ret []database.Digest
	conn := a.DB.Where("id IN (?)", IDs).
		Order(p.Order).Preload("Rule").Preload("Receipts").
		Find(&ret)
	if err := conn.Error; err != nil && !conn.RecordNotFound() {
		return nil, errors.Wrap(err, "finding digests")
	}

	return ret, nil
}
