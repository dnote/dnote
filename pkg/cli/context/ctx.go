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

// Package context defines dnote context
package context

import (
	"github.com/dnote/dnote/pkg/cli/database"
)

// DnoteCtx is a context holding the information of the current runtime
type DnoteCtx struct {
	HomeDir          string
	DnoteDir         string
	APIEndpoint      string
	Version          string
	DB               *database.DB
	SessionKey       string
	SessionKeyExpiry int64
	CipherKey        []byte
}
