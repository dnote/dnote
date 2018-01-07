package edit

import (
	"strconv"
	"time"

	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/log"
	"github.com/dnote-io/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var newContent string

var example = `
  * Edit the note by index in the current book
  dnote edit 3

  * Edit the note by index in a certain book
  dnote edit js 3

	* Skip the prompt by providing new content directly
	dntoe eidt js 3 -c "new content"`

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Edit a note or a book",
		Aliases: []string{"e"},
		Example: example,
		PreRunE: preRun,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&newContent, "content", "c", "", "The new content for the note")

	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing argument")
	}

	return nil
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dnote, err := core.GetDnote(ctx)
		if err != nil {
			return err
		}

		var targetBookName string
		var targetIdx int

		if len(args) == 1 {
			targetBookName, err = core.GetCurrentBook(ctx)
			if err != nil {
				return err
			}
			targetIdx, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
		} else if len(args) == 2 {
			targetBookName = args[0]
			targetIdx, err = strconv.Atoi(args[1])
			if err != nil {
				return err
			}
		}

		targetBook, exists := dnote[targetBookName]
		if !exists {
			return errors.Errorf("Book %s does not exist", targetBookName)
		}
		if targetIdx > len(targetBook.Notes)-1 {
			return errors.Errorf("Book %s does not have note with index %d", targetBookName, targetIdx)
		}
		targetNote := targetBook.Notes[targetIdx]

		if newContent == "" {
			log.Printf("content: %s\n", targetNote.Content)
			log.Printf("new content: ")

			newContent, err = utils.GetInput()
			if err != nil {
				return errors.Wrap(err, "Failed to get new content")
			}
		}

		ts := time.Now().Unix()

		targetNote.Content = utils.SanitizeContent(newContent)
		targetNote.EditedOn = ts
		targetBook.Notes[targetIdx] = targetNote
		dnote[targetBookName] = targetBook

		err = core.LogActionEditNote(ctx, targetNote.UUID, targetBook.Name, targetNote.Content, ts)
		if err != nil {
			return errors.Wrap(err, "Failed to log action")
		}

		err = core.WriteDnote(ctx, dnote)
		if err != nil {
			return errors.Wrap(err, "Failed to write dnote")
		}

		log.Info("edited the note")

		return nil
	}
}
