/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
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
