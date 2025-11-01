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
