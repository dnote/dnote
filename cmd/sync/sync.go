package sync

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/dnote/actions"
	"github.com/dnote/cli/core"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
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
	Actions  []actions.Action `json:"actions"`
	Bookmark int              `json:"bookmark"`
}

type syncPayload struct {
	Bookmark int    `json:"bookmark"`
	Actions  []byte `json:"actions"` // gziped
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		tx, err := ctx.DB.Begin()
		if err != nil {
			return errors.Wrap(err, "beginning a transaction")
		}

		config, err := core.ReadConfig(ctx)
		if err != nil {
			return errors.Wrap(err, "reading the config")
		}
		if config.APIKey == "" {
			log.Error("login required. please run `dnote login`\n")
			return nil
		}

		var bookmark int
		err = tx.QueryRow("SELECT value FROM system WHERE key = ?", "bookmark").Scan(&bookmark)
		if err != nil {
			return errors.Wrap(err, "getting bookmark")
		}

		actions, err := getLocalActions(tx)
		if err != nil {
			return errors.Wrap(err, "getting local actions")
		}

		payload, err := newPayload(actions, bookmark)
		if err != nil {
			return errors.Wrap(err, "getting the request payload")
		}

		log.Infof("writing changes (total %d).", len(actions))
		resp, err := postActions(ctx, config.APIKey, payload)
		if err != nil {
			return errors.Wrap(err, "posting to the server")
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "reading the response body")
		}

		if resp.StatusCode != http.StatusOK {
			bodyStr := string(body)

			fmt.Println("")
			return errors.Errorf("Server error: %s", bodyStr)
		}

		fmt.Println(" done.")

		var respData responseData
		err = json.Unmarshal(body, &respData)
		if err != nil {
			return errors.Wrap(err, "unmarshalling the payload")
		}

		log.Infof("resolving delta (total %d).", len(respData.Actions))
		err = core.ReduceAll(ctx, tx, respData.Actions)
		if err != nil {
			return errors.Wrap(err, "reducing returned actions")
		}
		fmt.Println(" done.")

		_, err = tx.Exec("UPDATE system SET value = ? WHERE key = ?", respData.Bookmark, "bookmark")
		if err != nil {
			return errors.Wrap(err, "updating the bookmark")
		}

		_, err = tx.Exec("DELETE FROM actions")
		if err != nil {
			return errors.Wrap(err, "clearing the action log")
		}

		tx.Commit()

		log.Success("success\n")

		if err := core.CheckUpdate(ctx); err != nil {
			log.Error(errors.Wrap(err, "automatically checking updates").Error())
		}

		return nil
	}
}

func newPayload(actions []actions.Action, bookmark int) (*bytes.Buffer, error) {
	compressedActions, err := compressActions(actions)
	if err != nil {
		return &bytes.Buffer{}, errors.Wrap(err, "compressing actions")
	}

	payload := syncPayload{
		Bookmark: bookmark,
		Actions:  compressedActions,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return &bytes.Buffer{}, errors.Wrap(err, "marshalling paylaod into JSON")
	}

	ret := bytes.NewBuffer(b)
	return ret, nil
}

func compressActions(actions []actions.Action) ([]byte, error) {
	b, err := json.Marshal(&actions)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling actions into JSON")
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)

	_, err = g.Write(b)
	if err != nil {
		return nil, errors.Wrap(err, "writing to gzip writer")
	}

	if err = g.Close(); err != nil {
		return nil, errors.Wrap(err, "closing gzip writer")
	}

	return buf.Bytes(), nil
}

func postActions(ctx infra.DnoteCtx, APIKey string, payload io.Reader) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/v1/sync", ctx.APIEndpoint)
	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return &http.Response{}, errors.Wrap(err, "forming an HTTP request")
	}

	req.Header.Set("Authorization", APIKey)
	req.Header.Set("CLI-Version", ctx.Version)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, errors.Wrap(err, "making a request")
	}

	return resp, nil
}

func getLocalActions(tx *sql.Tx) ([]actions.Action, error) {
	ret := []actions.Action{}

	rows, err := tx.Query("SELECT uuid, schema, type, data, timestamp FROM actions")
	if err != nil {
		return ret, errors.Wrap(err, "querying actions")
	}
	defer rows.Close()

	for rows.Next() {
		var action actions.Action

		err = rows.Scan(&action.UUID, &action.Schema, &action.Type, &action.Data, &action.Timestamp)
		if err != nil {
			return ret, errors.Wrap(err, "scanning a row")
		}

		ret = append(ret, action)
	}

	err = rows.Err()
	if err != nil {
		return ret, errors.Wrap(err, "scanning rows")
	}

	return ret, nil
}
