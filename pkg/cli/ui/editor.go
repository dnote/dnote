/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package ui provides the user interface for the program
package ui

import (
	"fmt"
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
		candidate := fmt.Sprintf("%s/%s", ctx.Paths.Cache, filename)

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

	b, err := os.ReadFile(fpath)
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
