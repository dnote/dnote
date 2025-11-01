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
