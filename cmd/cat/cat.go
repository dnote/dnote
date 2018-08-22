package cat

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
 * See the notes with index 2 from a book 'javascript'
 dnote cat javascript 2
 `

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("Incorrect number of arguments")
	}

	return nil
}

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cat <book name> <note index>",
		Aliases: []string{"c"},
		Short:   "See a note",
		Example: example,
		RunE:    newRun(ctx),
		PreRunE: preRun,
	}

	return cmd
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		dnote, err := core.GetDnote(ctx)
		if err != nil {
			return errors.Wrap(err, "reading dnote")
		}

		bookName := args[0]
		noteIdx, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.Wrapf(err, "parsing note index '%+v'", args[1])
		}

		book := dnote[bookName]
		note := book.Notes[noteIdx]

		log.Infof("book name: %s\n", bookName)
		log.Infof("note uuid: %s\n", note.UUID)
		log.Infof("created at: %s\n", time.Unix(note.AddedOn, 0).Format("Jan 2, 2006 3:04pm (MST)"))
		if note.EditedOn != 0 {
			log.Infof("updated at: %s\n", time.Unix(note.EditedOn, 0).Format("Jan 2, 2006 3:04pm (MST)"))
		}
		fmt.Printf("\n------------------------content------------------------\n")
		fmt.Printf("%s", note.Content)
		fmt.Printf("\n-------------------------------------------------------\n")

		return nil
	}
}
