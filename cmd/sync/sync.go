package sync

import (
	"bytes"
	"compress/gzip"
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
		config, err := core.ReadConfig(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to read the config")
		}
		timestamp, err := core.ReadTimestamp(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to read the timestamp")
		}
		actions, err := core.ReadActionLog(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to read the action log")
		}

		if config.APIKey == "" {
			log.Error("login required. please run `dnote login`\n")
			return nil
		}

		payload, err := getPayload(actions, timestamp)
		if err != nil {
			return errors.Wrap(err, "Failed to get dnote payload")
		}

		log.Infof("writing changes (total %d).", len(actions))
		resp, err := postActions(ctx, config.APIKey, payload)
		if err != nil {
			return errors.Wrap(err, "Failed to post to the server ")
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "Failed to read failed response body")
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
			return errors.Wrap(err, "Failed to unmarshal payload")
		}

		log.Infof("resolving delta (total %d).", len(respData.Actions))
		err = core.ReduceAll(ctx, respData.Actions)
		if err != nil {
			return errors.Wrap(err, "Failed to reduce returned actions")
		}
		fmt.Println(" done.")

		// Update bookmark
		ts, err := core.ReadTimestamp(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to read the timestamp")
		}
		ts.Bookmark = respData.Bookmark

		err = core.WriteTimestamp(ctx, ts)
		if err != nil {
			return errors.Wrap(err, "Failed to update bookmark")
		}

		log.Success("success\n")
		if err := core.ClearActionLog(ctx); err != nil {
			return errors.Wrap(err, "Failed to clear the action log")
		}

		if err := core.CheckUpdate(ctx); err != nil {
			log.Error(errors.Wrap(err, "automatically checking updates").Error())
		}

		return nil
	}
}

func getPayload(actions []actions.Action, timestamp infra.Timestamp) (*bytes.Buffer, error) {
	compressedActions, err := compressActions(actions)
	if err != nil {
		return &bytes.Buffer{}, errors.Wrap(err, "Failed to compress actions")
	}

	payload := syncPayload{
		Bookmark: timestamp.Bookmark,
		Actions:  compressedActions,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return &bytes.Buffer{}, errors.Wrap(err, "Failed to marshal paylaod into JSON")
	}

	ret := bytes.NewBuffer(b)
	return ret, nil
}

func compressActions(actions []actions.Action) ([]byte, error) {
	b, err := json.Marshal(&actions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal actions into JSON")
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)

	_, err = g.Write(b)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to write to gzip writer")
	}

	if err = g.Close(); err != nil {
		return nil, errors.Wrap(err, "Failed to close gzip writer")
	}

	return buf.Bytes(), nil
}

func postActions(ctx infra.DnoteCtx, APIKey string, payload io.Reader) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/v1/sync", ctx.APIEndpoint)
	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return &http.Response{}, errors.Wrap(err, "Failed to construct HTTP request")
	}

	req.Header.Set("Authorization", APIKey)
	req.Header.Set("CLI-Version", ctx.Version)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, errors.Wrap(err, "Failed to make request")
	}

	return resp, nil
}
