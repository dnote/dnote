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

// Package ui provides the user interface for the program
package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
)

// GetTmpContentPath returns the path to the temporary file containing
// content being added or edited
func GetTmpContentPath(ctx context.DnoteCtx) (string, error) {
	for i := 0; ; i++ {
		filename := fmt.Sprintf("%s_%d.%s", consts.TmpContentFileBase, i, consts.TmpContentFileExt)
		candidate := fmt.Sprintf("%s/%s", ctx.DnoteDir, filename)

		ok, err := utils.FileExists(candidate)
		if err != nil {
			return "", errors.Wrapf(err, "checking if file exists at %s", candidate)
		}
		if !ok {
			return candidate, nil
		}
	}
}

// getEditorCommand returns the system's editor command with appropriate flags,
// if necessary, to make the command wait until editor is close to exit.
func getEditorCommand() string {
	editor := os.Getenv("EDITOR")

	var ret string

	switch editor {
	case "atom":
		ret = "atom -w"
	case "subl":
		ret = "subl -n -w"
	case "mate":
		ret = "mate -w"
	case "vim":
		ret = "vim"
	case "nano":
		ret = "nano"
	case "emacs":
		ret = "emacs"
	case "nvim":
		ret = "nvim"
	default:
		ret = "vi"
	}

	return ret
}

func newEditorCmd(ctx context.DnoteCtx, fpath string) (*exec.Cmd, error) {
	args := strings.Fields(ctx.Editor)
	args = append(args, fpath)

	return exec.Command(args[0], args[1:]...), nil
}

// GetEditorInput gets the user input by launching a text editor and waiting for
// it to exit
func GetEditorInput(ctx context.DnoteCtx, fpath string) (string, error) {
	ok, err := utils.FileExists(fpath)
	if err != nil {
		return "", errors.Wrapf(err, "checking if the file exists at %s", fpath)
	}
	if !ok {
		f, err := os.Create(fpath)
		if err != nil {
			return "", errors.Wrap(err, "creating a temporary content file")
		}
		err = f.Close()
		if err != nil {
			return "", errors.Wrap(err, "closing the temporary content file")
		}
	}

	cmd, err := newEditorCmd(ctx, fpath)
	if err != nil {
		return "", errors.Wrap(err, "creating an editor command")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return "", errors.Wrapf(err, "launching an editor")
	}

	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, "waiting for the editor")
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return "", errors.Wrap(err, "reading the temporary content file")
	}

	err = os.Remove(fpath)
	if err != nil {
		return "", errors.Wrap(err, "removing the temporary content file")
	}

	raw := string(b)

	return raw, nil
}
