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

// Package logger provides an interface to transmit log messages to a system log service
package logger

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"log/syslog"
	"os"
)

var writer *syslog.Writer

// Init initilizes the syslog writer
func Init() error {
	var err error

	var endpoint string
	if os.Getenv("GO_ENV") == "PRODUCTION" {
		endpoint = "logs7.papertrailapp.com:37297"

		writer, err = syslog.Dial("udp", endpoint, syslog.LOG_DEBUG|syslog.LOG_KERN, "dnote-api")
		if err != nil {
			return errors.Wrap(err, "dialing syslog manager")
		}
	}

	return nil
}

// Info logs an info message
func Info(msg string, v ...interface{}) {
	m := fmt.Sprintf(msg, v...)
	log.Println(m)

	if writer == nil {
		return
	}

	if err := writer.Info(fmt.Sprintf("INFO: %s", m)); err != nil {
		log.Println(errors.Wrap(err, "transmiting log"))
	}
}

// Err logs an error message
func Err(msg string, v ...interface{}) {
	m := fmt.Sprintf(msg, v...)
	log.Println(m)

	if writer == nil {
		return
	}

	if err := writer.Err(fmt.Sprintf("ERROR: %s", m)); err != nil {
		log.Println(errors.Wrap(err, "transmiting log"))
	}
}

// Notice logs a notice message
func Notice(msg string, v ...interface{}) {
	m := fmt.Sprintf(msg, v...)
	log.Println(m)

	if writer == nil {
		return
	}

	if err := writer.Notice(fmt.Sprintf("NOTICE: %s", m)); err != nil {
		log.Println(errors.Wrap(err, "transmiting log"))
	}
}

// Debug logs a debug message
func Debug(msg string, v ...interface{}) {
	m := fmt.Sprintf(msg, v...)
	log.Println(m)

	if writer == nil {
		return
	}

	if err := writer.Debug(fmt.Sprintf("DEBUG: %s", m)); err != nil {
		log.Println(errors.Wrap(err, "transmiting log"))
	}
}
