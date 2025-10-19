/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
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

package assert

import (
	"bufio"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// WaitForPrompt waits for an expected prompt to appear in stdout with a timeout.
// Returns an error if the prompt is not found within the timeout period.
// Handles prompts with or without newlines by reading character by character.
func WaitForPrompt(stdout io.Reader, expectedPrompt string, timeout time.Duration) error {
	type result struct {
		found bool
		err   error
	}
	resultCh := make(chan result, 1)

	go func() {
		reader := bufio.NewReader(stdout)
		var buffer strings.Builder
		found := false

		for {
			b, err := reader.ReadByte()
			if err != nil {
				resultCh <- result{found: found, err: err}
				return
			}

			buffer.WriteByte(b)
			if strings.Contains(buffer.String(), expectedPrompt) {
				found = true
				break
			}
		}

		resultCh <- result{found: found, err: nil}
	}()

	select {
	case res := <-resultCh:
		if res.err != nil && res.err != io.EOF {
			return errors.Wrap(res.err, "reading stdout")
		}
		if !res.found {
			return errors.Errorf("expected prompt '%s' not found in stdout", expectedPrompt)
		}
		return nil
	case <-time.After(timeout):
		return errors.Errorf("timeout waiting for prompt '%s'", expectedPrompt)
	}
}

// RespondToPrompt is a helper that waits for a prompt and sends a response.
func RespondToPrompt(stdout io.Reader, stdin io.WriteCloser, expectedPrompt, response string, timeout time.Duration) error {
	if err := WaitForPrompt(stdout, expectedPrompt, timeout); err != nil {
		return err
	}

	if _, err := io.WriteString(stdin, response); err != nil {
		return errors.Wrap(err, "writing response to stdin")
	}

	return nil
}
