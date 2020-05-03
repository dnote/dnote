/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

// Package web provides handlers for the web application
package web

import (
	"net/http"

	"github.com/dnote/dnote/pkg/server/handlers"
	"github.com/dnote/dnote/pkg/server/tmpl"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	// ErrEmptyDatabase is an error for missing db in the context
	ErrEmptyDatabase = errors.New("No DB was provided")
	// ErrEmptyIndexHTML is an error for missing index.html content in the context
	ErrEmptyIndexHTML = errors.New("No index.html content was provided")
	// ErrEmptyRobotsTxt is an error for missing robots.txt content in the context
	ErrEmptyRobotsTxt = errors.New("No robots.txt content was provided")
	// ErrEmptyServiceWorkerJS is an error for missing service worker content in the context
	ErrEmptyServiceWorkerJS = errors.New("No service-worker.js content was provided")
	// ErrEmptyStaticFileSystem is an error for missing static filesystem in the context
	ErrEmptyStaticFileSystem = errors.New("No static filesystem was provided")
)

// Context contains contents of web assets
type Context struct {
	DB               *gorm.DB
	IndexHTML        []byte
	RobotsTxt        []byte
	ServiceWorkerJs  []byte
	StaticFileSystem http.FileSystem
}

// Handlers are a group of web handlers
type Handlers struct {
	GetRoot          http.HandlerFunc
	GetRobots        http.HandlerFunc
	GetServiceWorker http.HandlerFunc
	GetStatic        http.Handler
}

func validateContext(c Context) error {
	if c.DB == nil {
		return ErrEmptyDatabase
	}
	if c.IndexHTML == nil {
		return ErrEmptyIndexHTML
	}
	if c.RobotsTxt == nil {
		return ErrEmptyRobotsTxt
	}
	if c.ServiceWorkerJs == nil {
		return ErrEmptyServiceWorkerJS
	}
	if c.StaticFileSystem == nil {
		return ErrEmptyStaticFileSystem
	}

	return nil
}

// Init initializes the handlers
func Init(c Context) (Handlers, error) {
	if err := validateContext(c); err != nil {
		return Handlers{}, errors.Wrap(err, "validating context")
	}

	return Handlers{
		GetRoot:          getRootHandler(c),
		GetRobots:        getRobotsHandler(c),
		GetServiceWorker: getSWHandler(c),
		GetStatic:        getStaticHandler(c),
	}, nil
}

// getRootHandler returns an HTTP handler that serves the app shell
func getRootHandler(c Context) http.HandlerFunc {
	appShell, err := tmpl.NewAppShell(c.DB, c.IndexHTML)
	if err != nil {
		panic(errors.Wrap(err, "initializing app shell"))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// index.html must not be cached
		w.Header().Set("Cache-Control", "no-cache")

		buf, err := appShell.Execute(r)
		if err != nil {
			if errors.Cause(err) == tmpl.ErrNotFound {
				handlers.RespondNotFound(w)
			} else {
				handlers.DoError(w, "executing app shell", err, http.StatusInternalServerError)
			}
			return
		}

		w.Write(buf)
	}
}

// getRobotsHandler returns an HTTP handler that serves robots.txt
func getRobotsHandler(c Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(c.RobotsTxt)
	}
}

// getSWHandler returns an HTTP handler that serves service worker
func getSWHandler(c Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(c.ServiceWorkerJs)
	}
}

// getStaticHandler returns an HTTP handler that serves static files from a filesystem
func getStaticHandler(c Context) http.Handler {
	root := c.StaticFileSystem
	return http.StripPrefix("/static/", http.FileServer(root))
}
