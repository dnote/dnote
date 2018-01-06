package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const (
	letterRunes = "abcdefghipqrstuvwxyz0123456789"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateUID returns a uid
func GenerateUID() string {
	return uuid.NewV4().String()
}

func AskConfirmation(question string) (bool, error) {
	fmt.Printf("%s [Y/n]: ", question)

	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	ok := res == "y\n" || res == "Y\n" || res == "\n"

	return ok, nil
}

// FileExists checks if the file exists at the given path
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.Wrapf(err, "Failed to check if '%s' is directory", path)
	}

	return fileInfo.IsDir(), nil
}

// CopyFile copies a file from the src to dest
func CopyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "Failed to open the input file")
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return errors.Wrap(err, "Failed to create the output file")
	}

	if _, err = io.Copy(out, in); err != nil {
		return errors.Wrap(err, "Failed to copy the file content")
	}

	if err := out.Sync(); err != nil {
		return errors.Wrap(err, "Failed to flush the output file to disk")
	}

	fi, err := os.Stat(src)
	if err != nil {
		return errors.Wrap(err, "Failed to get file info for the input file")
	}

	if err := os.Chmod(dest, fi.Mode()); err != nil {
		return errors.Wrap(err, "Failed to copy permission to the output file")
	}

	// Close the output file
	if err := out.Close(); err != nil {
		return errors.Wrap(err, "Failed to close the output file")
	}

	return nil
}

// CopyDir copies a directory from src to dest, recursively copying nested
// directories
func CopyDir(src, dest string) error {
	srcPath := filepath.Clean(src)
	destPath := filepath.Clean(dest)

	fi, err := os.Stat(srcPath)
	if err != nil {
		return errors.Wrap(err, "Failed to get file info for the input")
	}
	if !fi.IsDir() {
		return errors.Wrap(err, "Source is not a directory")
	}

	_, err = os.Stat(dest)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "Failed to look up the destination")
	}

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return errors.Wrap(err, "Failed to create destination")
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.Wrap(err, "Failed to read directory listing for the input")
	}

	for _, entry := range entries {
		srcEntryPath := filepath.Join(srcPath, entry.Name())
		destEntryPath := filepath.Join(destPath, entry.Name())

		if entry.IsDir() {
			if err = CopyDir(srcEntryPath, destEntryPath); err != nil {
				return errors.Wrapf(err, "Failed to copy %s", entry.Name())
			}
		} else {
			if err = CopyFile(srcEntryPath, destEntryPath); err != nil {
				return errors.Wrapf(err, "Failed to copy %s", entry.Name())
			}
		}
	}

	return nil
}
