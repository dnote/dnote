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

// Package log provides interfaces to write structured logs
package log

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	fieldKeyLevel         = "level"
	fieldKeyMessage       = "msg"
	fieldKeyTimestamp     = "ts"
	fieldKeyUnixTimestamp = "ts_unix"

	// LevelDebug represents debug log level
	LevelDebug = "debug"
	// LevelInfo represents info log level
	LevelInfo = "info"
	// LevelWarn represents warn log level
	LevelWarn = "warn"
	// LevelError represents error log level
	LevelError = "error"
)

var (
	// currentLevel is the currently configured log level
	currentLevel = LevelInfo
)

// Fields represents a set of information to be included in the log
type Fields map[string]interface{}

// Entry represents a log entry
type Entry struct {
	Fields    Fields
	Timestamp time.Time
}

func newEntry(fields Fields) Entry {
	return Entry{
		Fields:    fields,
		Timestamp: time.Now().UTC(),
	}
}

// WithFields creates a log entry with the given fields
func WithFields(fields Fields) Entry {
	return newEntry(fields)
}

// SetLevel sets the global log level
func SetLevel(level string) {
	currentLevel = level
}

// GetLevel returns the current global log level
func GetLevel() string {
	return currentLevel
}

// shouldLog returns true if the given level should be logged based on currentLevel.
//
// Log level behavior (hierarchical):
//   - LevelDebug: shows all messages (debug, info, warn, error)
//   - LevelInfo: shows info, warn, and error messages
//   - LevelWarn: shows warn and error messages
//   - LevelError: shows only error messages
func shouldLog(level string) bool {
	// Debug level shows everything
	if currentLevel == LevelDebug {
		return true
	}

	// Info level shows info + warn + error
	if currentLevel == LevelInfo {
		return level == LevelInfo || level == LevelWarn || level == LevelError
	}

	// Warn level shows warn + error
	if currentLevel == LevelWarn {
		return level == LevelWarn || level == LevelError
	}

	// Error level shows only error
	if currentLevel == LevelError {
		return level == LevelError
	}

	return false
}

// Debug logs the given entry at a debug level
func (e Entry) Debug(msg string) {
	e.write(LevelDebug, msg)
}

// Info logs the given entry at an info level
func (e Entry) Info(msg string) {
	e.write(LevelInfo, msg)
}

// Warn logs the given entry at a warning level
func (e Entry) Warn(msg string) {
	e.write(LevelWarn, msg)
}

// Error logs the given entry at an error level
func (e Entry) Error(msg string) {
	e.write(LevelError, msg)
}

// ErrorWrap logs the given entry with the error message annotated by the given message
func (e Entry) ErrorWrap(err error, msg string) {
	m := fmt.Sprintf("%s: %v", msg, err)

	e.Error(m)
}

func (e Entry) formatJSON(level, msg string) []byte {
	data := make(Fields, len(e.Fields)+4)

	data[fieldKeyLevel] = level
	data[fieldKeyMessage] = msg
	data[fieldKeyTimestamp] = e.Timestamp
	data[fieldKeyUnixTimestamp] = e.Timestamp.Unix()

	for k, v := range e.Fields {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "formatting JSON: %v\n", err)
	}

	return serialized
}

func (e Entry) write(level, msg string) {
	if !shouldLog(level) {
		return
	}

	serialized := e.formatJSON(level, msg)

	_, err := fmt.Fprintln(os.Stderr, string(serialized))
	if err != nil {
		fmt.Fprintf(os.Stderr, "writing to stderr: %v\n", err)
	}
}

// Debug logs a debug message without additional fields
func Debug(msg string) {
	newEntry(Fields{}).Debug(msg)
}

// Info logs an info message without additional fields
func Info(msg string) {
	newEntry(Fields{}).Info(msg)
}

// Error logs an error message without additional fields
func Error(msg string) {
	newEntry(Fields{}).Error(msg)
}

// ErrorWrap logs an error message without additional fields. It annotates the given error's
// message with the given message
func ErrorWrap(err error, msg string) {
	newEntry(Fields{}).ErrorWrap(err, msg)
}
