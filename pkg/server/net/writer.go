/* Copyright (C) 2019, 2020, 2021, 2022 Monomax Software Pty Ltd
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

package net

import (
	"github.com/dnote/dnote/pkg/server/log"
	"net/http"
)

// LifecycleWriter  wraps http.ResponseWriter to track state of the http response.
// The optional interfaces of http.ResponseWriter are lost because of the wrapping, and
// such interfaces should be implemented if needed. (i.e. http.Pusher, http.Flusher, etc.)
type LifecycleWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader wraps the WriteHeader call and marks the response state as done.
func (w *LifecycleWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// IsHeaderWritten returns true if a response has been written.
func IsHeaderWritten(w http.ResponseWriter) bool {
	if lw, ok := w.(*LifecycleWriter); ok {
		return lw.StatusCode != 0
	}

	// the response writer must have been wrapped in the middleware chain.
	log.Error("unable to log because writer is not a LifecycleWriter")
	return false
}
