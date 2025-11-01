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

// Package output provides functions to print informations on the terminal
// in a consistent manner
package output

import (
	"fmt"
	"time"

	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/log"
)

// NoteInfo prints a note information
func NoteInfo(info database.NoteInfo) {
	log.Infof("book name: %s\n", info.BookLabel)
	log.Infof("created at: %s\n", time.Unix(0, info.AddedOn).Format("Jan 2, 2006 3:04pm (MST)"))
	if info.EditedOn != 0 {
		log.Infof("updated at: %s\n", time.Unix(0, info.EditedOn).Format("Jan 2, 2006 3:04pm (MST)"))
	}
	log.Infof("note id: %d\n", info.RowID)
	log.Infof("note uuid: %s\n", info.UUID)

	fmt.Printf("\n------------------------content------------------------\n")
	fmt.Printf("%s", info.Content)
	fmt.Printf("\n-------------------------------------------------------\n")
}

func NoteContent(info database.NoteInfo) {
	fmt.Printf("%s", info.Content)
}

// BookInfo prints a note information
func BookInfo(info database.BookInfo) {
	log.Infof("book name: %s\n", info.Name)
	log.Infof("book id: %d\n", info.RowID)
	log.Infof("book uuid: %s\n", info.UUID)
}
