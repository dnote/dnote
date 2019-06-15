/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package context

import (
	"encoding/base64"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/pkg/errors"
)

// GetCipherKey retrieves the cipher key and decode the base64 into bytes.
func (ctx DnoteCtx) GetCipherKey() ([]byte, error) {
	tx, err := ctx.DB.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "beginning transaction")
	}

	var cipherKeyB64 string
	err = database.GetSystem(tx, consts.SystemCipherKey, &cipherKeyB64)
	if err != nil {
		return []byte{}, errors.Wrap(err, "getting enc key")
	}

	cipherKey, err := base64.StdEncoding.DecodeString(cipherKeyB64)
	if err != nil {
		return nil, errors.Wrap(err, "decoding cipherKey from base64")
	}

	return cipherKey, nil
}
