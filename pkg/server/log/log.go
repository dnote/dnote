/* Copyright (C) 2019 Monomax Software Pty Ltd
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

	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
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

// Info logs the given entry at an info level
func (e Entry) Info(msg string) {
	e.write(levelInfo, msg)
}

// Warn logs the given entry at a warning level
func (e Entry) Warn(msg string) {
	e.write(levelWarn, msg)
}

// Error logs the given entry at an error level
func (e Entry) Error(msg string) {
	e.write(levelError, msg)
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
	serialized := e.formatJSON(level, msg)

	_, err := fmt.Fprintln(os.Stderr, string(serialized))
	if err != nil {
		fmt.Fprintf(os.Stderr, "writing to stderr: %v\n", err)
	}
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
