package sync

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dnote-io/cli/core"
	"github.com/dnote-io/cli/infra"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
  dnote sync`

func NewCmd(ctx infra.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Sync dnote with the dnote server",
		Example: example,
		RunE:    newRun(ctx),
	}

	return cmd
}

type responseData struct {
	Actions  []core.Action `json:"actions"`
	Bookmark int           `json:"bookmark"`
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

		if config.APIKey == "" {
			fmt.Println("Login required. Please run `dnote login`")
			return nil
		}

		fmt.Println("Compressing dnote...")
		payload, err := getPayload(ctx, timestamp)
		if err != nil {
			return errors.Wrap(err, "Failed to get dnote payload")
		}

		fmt.Println("Syncing with the server...")
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

			fmt.Printf("Failed to sync on the server: %s", bodyStr)
			return errors.New(bodyStr)
		}

		var respData responseData
		err = json.Unmarshal(body, &respData)
		if err != nil {
			return errors.Wrap(err, "Failed to unmarshal payload")
		}

		err = core.ReduceAll(ctx, respData.Actions)
		if err != nil {
			return errors.Wrap(err, "Failed to reduce returned actions")
		}

		// Update bookmark
		ts, err := core.ReadTimestamp(ctx)
		ts.Bookmark = respData.Bookmark

		err = core.WriteTimestamp(ctx, ts)
		if err != nil {
			return errors.Wrap(err, "Failed to update bookmark")
		}

		fmt.Println("Successfully synced all notes")
		if err := core.ClearActionLog(ctx); err != nil {
			return errors.Wrap(err, "Failed to clear the action log")
		}

		return nil
	}
}

func getPayload(ctx infra.DnoteCtx, timestamp infra.Timestamp) (*bytes.Buffer, error) {
	actions, err := compressActions(ctx)
	if err != nil {
		return &bytes.Buffer{}, errors.Wrap(err, "Failed to compress actions")
	}

	payload := syncPayload{
		Bookmark: timestamp.Bookmark,
		Actions:  actions,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return &bytes.Buffer{}, errors.Wrap(err, "Failed to marshal paylaod into JSON")
	}

	ret := bytes.NewBuffer(b)
	return ret, nil
}

func compressActions(ctx infra.DnoteCtx) ([]byte, error) {
	b, err := core.ReadActionLogContent(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read the action log content")
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

func postActions(ctx infra.DnoteCtx, APIKey string, payload *bytes.Buffer) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/v1/sync", ctx.APIEndpoint)
	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return &http.Response{}, errors.Wrap(err, "Failed to construct HTTP request")
	}

	req.Header.Set("Authorization", APIKey)
	req.Header.Set("CLI-Version", core.Version)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, errors.Wrap(err, "Failed to make request")
	}

	return resp, nil
}
