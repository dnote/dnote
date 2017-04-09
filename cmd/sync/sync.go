package sync

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dnote-io/cli/utils"
)

func compressDnote() (*bytes.Buffer, error) {
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
	config, err := utils.ReadConfig()
	if err != nil {
		return err
	}

	if config.APIKey == "" {
		fmt.Println("Login required. Please run `dnote login`")
		return nil
	}

	fmt.Println("Compressing dnote...")
	payload, err := compressDnote()
	if err != nil {
		return err
	}

	fmt.Println("Syncing...")
	req, err := http.NewRequest("POST", "http://127.0.0.1:3030/sync", payload)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", config.APIKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		bodyStr := string(body)

		fmt.Printf("Failed to sync: %s", bodyStr)
	} else {
		fmt.Println("Successfully synced all notes")
	}

	return nil
}
