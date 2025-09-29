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

package assets

import (
	"embed"
	"github.com/pkg/errors"
	"io/fs"
)

//go:embed static
var staticFiles embed.FS

// GetStaticFS returns a filesystem for static files, with
// all files situated in the root of the filesystem
func GetStaticFS() (fs.FS, error) {
	subFs, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, errors.Wrap(err, "getting sub filesystem")
	}

	return subFs, nil
}

// MustGetHTTP500ErrorPage returns the content of HTML file for HTTP 500 error
func MustGetHTTP500ErrorPage() []byte {
	ret, err := staticFiles.ReadFile("static/500.html")
	if err != nil {
		panic(errors.Wrap(err, "reading HTML file for 500 HTTP error"))
	}

	return ret
}
