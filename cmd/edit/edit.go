package edit

import (
	"io/ioutil"
	"strconv"
	"time"

	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/dnote-io/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var newContent string

var example = `
  * Edit the note by index in a book
  dnote edit js 3

	* Skip the prompt by providing new content directly
	dntoe edit js 3 -c "new content"`

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
	if len(args) != 2 {
		return errors.New("Incorrect number of argument")
	}

	return nil
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dnote, err := core.GetDnote(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to read dnote")
		}

		targetBookName := args[0]
		targetIdx, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.Wrapf(err, "Failed to parse the given index %+v", args[1])
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
			fpath := core.GetDnoteTmpContentPath(ctx)

			e := ioutil.WriteFile(fpath, []byte(targetNote.Content), 0644)
			if err != nil {
				return errors.Wrap(err, "Failed to prepare editor content")
			}

			e = core.GetEditorInput(ctx, fpath, &newContent)
			if e != nil {
				return errors.Wrap(err, "Failed to get editor input")
			}

		}

		if targetNote.Content == newContent {
			return errors.New("Nothing changed")
		}

		ts := time.Now().Unix()

		targetNote.Content = core.SanitizeContent(newContent)
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

		log.Printf("new content: %s\n", newContent)
		log.Success("edited the note\n")

		return nil
	}
}
