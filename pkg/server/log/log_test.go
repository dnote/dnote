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

package log

import (
	"testing"
)

func TestSetLevel(t *testing.T) {
	// Reset to default after test
	defer SetLevel(LevelInfo)

	SetLevel(LevelDebug)
	if currentLevel != LevelDebug {
		t.Errorf("Expected level %s, got %s", LevelDebug, currentLevel)
	}

	SetLevel(LevelError)
	if currentLevel != LevelError {
		t.Errorf("Expected level %s, got %s", LevelError, currentLevel)
	}
}

func TestShouldLog(t *testing.T) {
	// Reset to default after test
	defer SetLevel(LevelInfo)

	testCases := []struct {
		currentLevel string
		logLevel     string
		expected     bool
		description  string
	}{
		// Debug level shows everything
		{LevelDebug, LevelDebug, true, "debug level should show debug"},
		{LevelDebug, LevelInfo, true, "debug level should show info"},
		{LevelDebug, LevelWarn, true, "debug level should show warn"},
		{LevelDebug, LevelError, true, "debug level should show error"},

		// Info level shows info + warn + error
		{LevelInfo, LevelDebug, false, "info level should not show debug"},
		{LevelInfo, LevelInfo, true, "info level should show info"},
		{LevelInfo, LevelWarn, true, "info level should show warn"},
		{LevelInfo, LevelError, true, "info level should show error"},

		// Warn level shows warn + error
		{LevelWarn, LevelDebug, false, "warn level should not show debug"},
		{LevelWarn, LevelInfo, false, "warn level should not show info"},
		{LevelWarn, LevelWarn, true, "warn level should show warn"},
		{LevelWarn, LevelError, true, "warn level should show error"},

		// Error level shows only error
		{LevelError, LevelDebug, false, "error level should not show debug"},
		{LevelError, LevelInfo, false, "error level should not show info"},
		{LevelError, LevelWarn, false, "error level should not show warn"},
		{LevelError, LevelError, true, "error level should show error"},
	}

	for _, tc := range testCases {
		SetLevel(tc.currentLevel)
		result := shouldLog(tc.logLevel)
		if result != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.description, tc.expected, result)
		}
	}
}
