/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// ReadFileAbs reads the content of the file with the given file path by resolving
// it as an absolute path
func ReadFileAbs(relpath string) []byte {
	fp, err := filepath.Abs(relpath)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	return b
}

// FileExists checks if the file exists at the given path
func FileExists(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, errors.Wrap(err, "getting file info")
}

// CopyDir copies a directory from src to dest, recursively copying nested
// directories
func CopyDir(src, dest string) error {
	srcPath := filepath.Clean(src)
	destPath := filepath.Clean(dest)

	fi, err := os.Stat(srcPath)
	if err != nil {
		return errors.Wrap(err, "getting the file info for the input")
	}
	if !fi.IsDir() {
		return errors.New("source is not a directory")
	}

	_, err = os.Stat(dest)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "looking up the destination")
	}

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return errors.Wrap(err, "creating destination")
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.Wrap(err, "reading the directory listing for the input")
	}

	for _, entry := range entries {
		srcEntryPath := filepath.Join(srcPath, entry.Name())
		destEntryPath := filepath.Join(destPath, entry.Name())

		if entry.IsDir() {
			if err = CopyDir(srcEntryPath, destEntryPath); err != nil {
				return errors.Wrapf(err, "copying %s", entry.Name())
			}
		} else {
			if err = CopyFile(srcEntryPath, destEntryPath); err != nil {
				return errors.Wrapf(err, "copying %s", entry.Name())
			}
		}
	}

	return nil
}

// CopyFile copies a file from the src to dest
func CopyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "opening the input file")
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return errors.Wrap(err, "creating the output file")
	}

	if _, err = io.Copy(out, in); err != nil {
		return errors.Wrap(err, "copying the file content")
	}

	if err = out.Sync(); err != nil {
		return errors.Wrap(err, "flushing the output file to disk")
	}

	fi, err := os.Stat(src)
	if err != nil {
		return errors.Wrap(err, "getting the file info for the input file")
	}

	if err = os.Chmod(dest, fi.Mode()); err != nil {
		return errors.Wrap(err, "copying permission to the output file")
	}

	// Close the output file
	if err = out.Close(); err != nil {
		return errors.Wrap(err, "closing the output file")
	}

	return nil
}
