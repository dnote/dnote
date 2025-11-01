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

package edit

import (
	"strings"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/output"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/dnote/dnote/pkg/cli/validate"
	"github.com/pkg/errors"
)

func validateRunBookFlags() error {
	if contentFlag != "" {
		return errors.New("--content is invalid for editing a book")
	}
	if bookFlag != "" {
		return errors.New("--book is invalid for editing a book")
	}

	return nil
}

func waitEditorBookName(ctx context.DnoteCtx) (string, error) {
	fpath, err := ui.GetTmpContentPath(ctx)
	if err != nil {
		return "", errors.Wrap(err, "getting temporarily content file path")
	}

	c, err := ui.GetEditorInput(ctx, fpath)
	if err != nil {
		return "", errors.Wrap(err, "getting editor input")
	}

	// remove the newline at the end because files end with linebreaks in POSIX
	c = strings.TrimSuffix(c, "\n")
	c = strings.TrimSuffix(c, "\r\n")

	return c, nil
}

func getName(ctx context.DnoteCtx) (string, error) {
	if nameFlag != "" {
		return nameFlag, nil
	}

	c, err := waitEditorBookName(ctx)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get editor input")
	}

	return c, nil
}

func runBook(ctx context.DnoteCtx, bookName string) error {
	err := validateRunBookFlags()
	if err != nil {
		return errors.Wrap(err, "validating flags.")
	}

	db := ctx.DB
	uuid, err := database.GetBookUUID(db, bookName)
	if err != nil {
		return errors.Wrap(err, "getting book uuid")
	}

	name, err := getName(ctx)
	if err != nil {
		return errors.Wrap(err, "getting name")
	}

	err = validate.BookName(name)
	if err != nil {
		return errors.Wrap(err, "validating book name")
	}

	tx, err := ctx.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	err = database.UpdateBookName(tx, uuid, name)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "updating the book name")
	}

	bookInfo, err := database.GetBookInfo(tx, uuid)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "getting book info")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "committing a transaction")
	}

	log.Success("edited the book\n")
	output.BookInfo(bookInfo)

	return nil
}
