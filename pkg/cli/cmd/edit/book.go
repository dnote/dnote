/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
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
