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

package presenters

import (
	"time"

	"github.com/dnote/dnote/pkg/server/database"
)

// Book is a result of PresentBooks
type Book struct {
	UUID      string    `json:"uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Label     string    `json:"label"`
}

// PresentBook presents a book
func PresentBook(book database.Book) Book {
	return Book{
		UUID:      book.UUID,
		USN:       book.USN,
		CreatedAt: FormatTS(book.CreatedAt),
		UpdatedAt: FormatTS(book.UpdatedAt),
		Label:     book.Label,
	}
}

// PresentBooks presents books
func PresentBooks(books []database.Book) []Book {
	ret := []Book{}

	for _, book := range books {
		p := PresentBook(book)
		ret = append(ret, p)
	}

	return ret
}
