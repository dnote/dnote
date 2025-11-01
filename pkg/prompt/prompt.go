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

// Package prompt provides utilities for interactive yes/no prompts
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// FormatQuestion formats a yes/no question with the appropriate choice indicator
func FormatQuestion(question string, optimistic bool) string {
	choices := "(y/N)"
	if optimistic {
		choices = "(Y/n)"
	}
	return fmt.Sprintf("%s %s", question, choices)
}

// ReadYesNo reads and parses a yes/no response from the given reader.
// Returns true if confirmed, respecting optimistic mode.
// In optimistic mode, empty input is treated as confirmation.
func ReadYesNo(r io.Reader, optimistic bool) (bool, error) {
	reader := bufio.NewReader(r)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	confirmed := input == "y"

	if optimistic {
		confirmed = confirmed || input == ""
	}

	return confirmed, nil
}
