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
