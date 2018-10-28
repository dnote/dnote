package sync

import (
	"database/sql"
	//"io"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/dnote/cli/migrate"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
  dnote sync`

// NewCmd returns a new sync command
func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Aliases: []string{"s"},
		Short:   "Sync dnote with the dnote server",
		Example: example,
		RunE:    newRun(ctx),
	}

	return cmd
}

type responseData struct {
	Bookmark int `json:"bookmark"`
}

type syncPayload struct {
	Bookmark int `json:"bookmark"`
}

func getLastSyncAt(ctx infra.DnoteCtx) (int, error) {
	ret := 0

	db := ctx.DB

	var count int
	err := db.QueryRow("SELECT count(*) FROM system WHERE key = ?", infra.SystemLastSyncAt).Scan(&count)
	if err != nil {
		return ret, errors.Wrap(err, "counting last sync time")
	}

	if count == 0 {
		return ret, nil
	}

	err = db.QueryRow("SELECT value FROM system WHERE key = ?", infra.SystemLastSyncAt).Scan(&ret)
	if err != nil {
		return ret, errors.Wrap(err, "querying last sync time")
	}

	return ret, nil
}

func getLastMaxUSN(ctx infra.DnoteCtx) (int, error) {
	ret := 0

	db := ctx.DB

	var count int
	err := db.QueryRow("SELECT count(*) FROM system WHERE key = ?", infra.SystemLastMaxUSN).Scan(&count)
	if err != nil {
		return ret, errors.Wrap(err, "counting last user max_usn")
	}

	if count == 0 {
		return ret, nil
	}

	err = db.QueryRow("SELECT value FROM system WHERE key = ?", infra.SystemLastMaxUSN).Scan(&ret)
	if err != nil {
		return ret, errors.Wrap(err, "querying last user max_usn")
	}

	return ret, nil
}

type syncStateResp struct {
	FullSyncBefore int `json:"full_sync_before"`
	MaxUSN         int `json:"max_usn"`
}

func getSyncState(apiKey string, ctx infra.DnoteCtx) (syncStateResp, error) {
	var ret syncStateResp

	res, err := utils.DoAuthorizedReq(ctx, apiKey, "GET", "/v1/sync/state", "")
	if err != nil {
		return ret, errors.Wrap(err, "constructing http request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ret, errors.Wrap(err, "reading the response body")
	}

	if err = json.Unmarshal(body, &ret); err != nil {
		return ret, errors.Wrap(err, "unmarshalling the payload")
	}

	return ret, nil
}

// syncFragNote represents a note in a sync fragment and contains only the necessary information
// for the client to sync the note locally
type syncFragNote struct {
	UUID      string    `json:"uuid"`
	BookUUID  string    `json:"book_uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AddedOn   int64     `json:"added_on"`
	EditedOn  int64     `json:"edited_on"`
	Content   string    `json:"content"`
	Public    bool      `json:"public"`
	Deleted   bool      `json:"deleted"`
}

// syncFragBook represents a book in a sync fragment and contains only the necessary information
// for the client to sync the note locally
type syncFragBook struct {
	UUID      string    `json:"uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AddedOn   int64     `json:"added_on"`
	Label     string    `json:"label"`
	Deleted   bool      `json:"deleted"`
}

// syncFragment contains a piece of information about the server's state.
type syncFragment struct {
	FragMaxUSN  int            `json:"frag_max_usn"`
	UserMaxUSN  int            `json:"user_max_usn"`
	CurrentTime int64          `json:"current_time"`
	Notes       []syncFragNote `json:"notes"`
	Books       []syncFragBook `json:"books"`
}

// syncList is an aggregation of resources represented in the sync fragments
type syncList struct {
	Notes []syncFragNote
	Books []syncFragBook
}

type getSyncFragmentResp struct {
	Fragment syncFragment `json:"fragment"`
}

func getSyncFragments(ctx infra.DnoteCtx, apiKey string, afterUSN int) ([]syncFragment, error) {
	var buf []syncFragment

	nextAfterUSN := afterUSN

	for {
		v := url.Values{}
		v.Set("after_usn", strconv.Itoa(nextAfterUSN))
		queryStr := v.Encode()

		path := fmt.Sprintf("/v1/sync/fragment?%s", queryStr)
		res, err := utils.DoAuthorizedReq(ctx, apiKey, "GET", path, "")

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return buf, errors.Wrap(err, "reading the response body")
		}

		var resp getSyncFragmentResp
		if err = json.Unmarshal(body, &resp); err != nil {
			return buf, errors.Wrap(err, "unmarshalling the payload")
		}

		frag := resp.Fragment
		buf = append(buf, frag)

		nextAfterUSN = frag.FragMaxUSN

		// if there is no more data, break
		if nextAfterUSN == 0 {
			break
		}
	}

	return buf, nil
}

func getSyncList(fragments []syncFragment) syncList {
	var notes []syncFragNote
	var books []syncFragBook

	for _, fragment := range fragments {
		notes = append(notes, fragment.Notes...)
		books = append(books, fragment.Books...)
	}

	return syncList{
		Notes: notes,
		Books: books,
	}
}

func resolveNotes(tx *sql.Tx, fragNotes []syncFragNote) error {
	for _, n := range fragNotes {
		var localUSN int
		var dirty bool
		err := tx.QueryRow("SELECT usn, dirty FROM notes WHERE uuid = ?", n.UUID).Scan(&localUSN, &dirty)
		if err != nil && err != sql.ErrNoRows {
			return errors.Wrapf(err, "getting local note %s", n.UUID)
		}

		// if note exists in the server and does not exist in the client, insert the note.
		if err == sql.ErrNoRows {
			if _, err := tx.Exec("INSERT INTO notes (uuid, content, usn, added_on) VALUES (?, ?, ?, ?, ?)",
				n.UUID, n.Content, n.USN, n.AddedOn); err != nil {
				return errors.Wrapf(err, "inserting note with uuid %s", n.UUID)
			}
		} else if n.USN > localUSN {
			// automatically resolve conflicts by letting server overwrite local changes
			// IDEA: if needed, report conflict instead of overwriting
			if _, err := tx.Exec("UPDATE notes SET usn = ?, book_uuid = ?, content = ?, edited_on = ?, deleted = ? WHERE uuid = ?",
				localUSN, n.BookUUID, n.Content, n.EditedOn, n.Deleted, n.UUID); err != nil {
				return errors.Wrapf(err, "updating local note %s", n.UUID)
			}
		}
	}

	return nil
}

func fullSync(ctx infra.DnoteCtx, apiKey string, afterUSN int) error {
	fragments, err := getSyncFragments(ctx, apiKey, afterUSN)
	if err != nil {
		return errors.Wrap(err, "getting sync fragments")
	}

	list := getSyncList(fragments)

	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	if err := resolveNotes(tx, list.Notes); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "resolving notes")
	}

	tx.Commit()

	return nil
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		config, err := core.ReadConfig(ctx)
		if err != nil {
			return errors.Wrap(err, "reading the config")
		}
		if config.APIKey == "" {
			log.Error("login required. please run `dnote login`\n")
			return nil
		}

		if err := migrate.Run(ctx, migrate.RemoteSequence, migrate.RemoteMode); err != nil {
			return errors.Wrap(err, "running remote migrations")
		}

		syncState, err := getSyncState(config.APIKey, ctx)
		if err != nil {
			return errors.Wrap(err, "getting the sync state from the server")
		}
		lastSyncAt, err := getLastSyncAt(ctx)
		if err != nil {
			return errors.Wrap(err, "getting the last sync time")
		}
		lastMaxUSN, err := getLastMaxUSN(ctx)
		if err != nil {
			return errors.Wrap(err, "getting the last max_usn")
		}

		if lastSyncAt < syncState.FullSyncBefore {
			// full sync
		} else if lastMaxUSN == syncState.MaxUSN {
			// skip to send changes
		} else {
			// incremental sync
		}

		log.Success("success\n")

		if err := core.CheckUpdate(ctx); err != nil {
			log.Error(errors.Wrap(err, "automatically checking updates").Error())
		}

		return nil
	}
}
