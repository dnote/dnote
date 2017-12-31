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
	Actions []core.Action `json:"actions"`
}

func newRun(ctx infra.DnoteCtx) core.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		config, err := core.ReadConfig(ctx)
		if err != nil {
			return err
		}

		if config.APIKey == "" {
			fmt.Println("Login required. Please run `dnote login`")
			return nil
		}

		fmt.Println("Compressing dnote...")
		payload, err := compressActions(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to compress dnote")
		}

		fmt.Println("Syncing...")

		endpoint := fmt.Sprintf("%s/v1/sync", ctx.APIEndpoint)
		req, err := http.NewRequest("POST", endpoint, payload)
		if err != nil {
			return errors.Wrap(err, "Failed to construct HTTP request")
		}
		req.Header.Set("Authorization", config.APIKey)
		req.Header.Set("CLI-Version", core.Version)

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrap(err, "Failed to make request")
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

		fmt.Println("resp body", string(body))

		var respData responseData
		if err := json.Unmarshal(body, &respData); err != nil {
			fmt.Println(err.Error())
			return errors.Wrap(err, "Failed to unmarshal payload")
		}

		// TODO: transaction
		for _, action := range respData.Actions {
			if err := core.Reduce(ctx, action); err != nil {
				return errors.Wrap(err, "Failed to reduce action")
			}
		}

		fmt.Println("Successfully synced all notes")
		if err := core.ClearActionLog(ctx); err != nil {
			return errors.Wrap(err, "Failed to clear the action log")
		}

		return nil
	}
}

func compressActions(ctx infra.DnoteCtx) (*bytes.Buffer, error) {
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

	return &buf, nil
}
