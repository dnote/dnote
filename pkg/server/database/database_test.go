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

package database

import (
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/log"
	"gorm.io/gorm/logger"
)

func TestGetDBLogLevel(t *testing.T) {
	testCases := []struct {
		name     string
		level    string
		expected logger.LogLevel
	}{
		{
			name:     "debug level maps to Info",
			level:    log.LevelDebug,
			expected: logger.Info,
		},
		{
			name:     "info level maps to Silent",
			level:    log.LevelInfo,
			expected: logger.Silent,
		},
		{
			name:     "warn level maps to Warn",
			level:    log.LevelWarn,
			expected: logger.Warn,
		},
		{
			name:     "error level maps to Error",
			level:    log.LevelError,
			expected: logger.Error,
		},
		{
			name:     "unknown level maps to Silent",
			level:    "unknown",
			expected: logger.Silent,
		},
		{
			name:     "empty string maps to Silent",
			level:    "",
			expected: logger.Silent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getDBLogLevel(tc.level)
			assert.Equal(t, result, tc.expected, "log level mismatch")
		})
	}
}
