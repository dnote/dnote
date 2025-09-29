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

package controllers

import (
	"net/http"
	"strings"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/views"
)

// NewStatic creates a new Static controller.
func NewStatic(app *app.App, viewEngine *views.Engine) *Static {
	return &Static{
		NotFoundView: viewEngine.NewView(app, views.Config{Title: "Not Found", Layout: "base"}, "static/not_found"),
	}
}

// Static is a static controller
type Static struct {
	NotFoundView *views.View
}

// NotFound is a catch-all handler for requests with no matching handler
func (s *Static) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "text/html") {
		s.NotFoundView.Render(w, r, nil, http.StatusOK)
	} else {
		statusText := http.StatusText(http.StatusNotFound)
		w.Write([]byte(statusText))
	}
}
