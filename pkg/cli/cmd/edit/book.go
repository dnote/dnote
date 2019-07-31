package edit

import (
	"strings"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/output"
	"github.com/dnote/dnote/pkg/cli/ui"
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

func waitEditorBookName(ctx context.DnoteCtx) error {
	fpath, err := ui.GetTmpContentPath(ctx)
	if err != nil {
		return errors.Wrap(err, "getting temporarily content file path")
	}

	if err := ui.GetEditorInput(ctx, fpath, &nameFlag); err != nil {
		return errors.Wrap(err, "getting editor input")
	}

	// remove the newline at the end because files end with linebreaks in POSIX
	nameFlag = strings.TrimSuffix(nameFlag, "\n")
	nameFlag = strings.TrimSuffix(nameFlag, "\r\n")

	return nil
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

	if nameFlag == "" {
		err := waitEditorBookName(ctx)
		if err != nil {
			return errors.Wrap(err, "getting content from editor")
		}
	}

	tx, err := ctx.DB.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	err = database.UpdateBookName(tx, uuid, nameFlag)
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
