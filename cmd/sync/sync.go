package sync

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/dnote-io/cli/utils"
)

func getRequestPayload() (*bytes.Buffer, error) {
	b, err := utils.ReadNoteContent()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)

	if _, err := g.Write(b); err != nil {
		return nil, err
	}

	if err = g.Close(); err != nil {
		return nil, err
	}

	return &buf, nil
}

func Sync() error {
	payload, err := getRequestPayload()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:3030/sync", payload)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}
