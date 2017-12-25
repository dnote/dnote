package sync

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"

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

func newRun(ctx infra.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		config, err := infra.ReadConfig(ctx)
		if err != nil {
			return err
		}

		if config.APIKey == "" {
			fmt.Println("Login required. Please run `dnote login`")
			return nil
		}

		fmt.Println("Compressing dnote...")
		payload, err := compressDnote(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to compress dnote")
		}

		fmt.Println("Syncing...")
		//endpoint := "http://api.dnote.io/v1/sync"
		endpoint := "http://127.0.0.1:3030/v1/sync"
		req, err := http.NewRequest("POST", endpoint, payload)
		if err != nil {
			return errors.Wrap(err, "Failed to construct HTTP request")
		}
		req.Header.Set("Authorization", config.APIKey)
		req.Header.Set("CLI-Version", infra.Version)

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrap(err, "Failed to make request")
		}

		if resp.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, "Failed to read failed response body")
			}
			bodyStr := string(body)

			fmt.Printf("Failed to sync: %s", bodyStr)
			return errors.New(bodyStr)
		}

		fmt.Println("Successfully synced all notes")

		return nil
	}

}

func compressDnote(ctx infra.DnoteCtx) (*bytes.Buffer, error) {
	b, err := infra.ReadNoteContent(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read note content")
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)

	_, err = g.Write(b)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to write note to gzip writer")
	}

	if err = g.Close(); err != nil {
		return nil, errors.Wrap(err, "Failed to close gzip writer")
	}

	return &buf, nil
}
