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

package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/buildinfo"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/gorilla/csrf"
)

const (
	// TemplateExt is the template extension
	TemplateExt string = ".gohtml"
)

const (
	siteTitle = "Dnote"
)

// Config is a view config
type Config struct {
	Title          string
	Layout         string
	HeaderTemplate string
	HelperFuncs    map[string]interface{}
	AlertInBody    bool
	Clock          clock.Clock
}

type viewCtx struct {
	Clock  clock.Clock
	Config Config
}

func newViewCtx(c Config) viewCtx {
	return viewCtx{
		Clock:  c.getClock(),
		Config: c,
	}
}

func (c Config) getLayout() string {
	if c.Layout == "" {
		return "base"
	}

	return c.Layout
}

func (c Config) getClock() clock.Clock {
	if c.Clock != nil {
		return c.Clock
	}

	return clock.New()
}

// View holds the information about a view
type View struct {
	Template *template.Template
	Layout   string
	// AlertInBody specifies if alert should be set in the body instead of the header
	AlertInBody bool
	App         *app.App
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil, http.StatusOK)
}

// Render is used to render the view with the predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data *Data, statusCode int) {
	w.Header().Set("Content-Type", "text/html")

	var vd Data
	if data != nil {
		vd = *data
	}

	if alert := getAlert(r); alert != nil {
		vd.PutAlert(*alert, v.AlertInBody)
		clearAlert(w)
	}

	vd.User = context.User(r.Context())

	// Put user data in Yield
	if vd.Yield == nil {
		vd.Yield = map[string]interface{}{}
	}
	if vd.User != nil {
		vd.Yield["Email"] = vd.User.Email.String
	}
	vd.Yield["CurrentPath"] = r.URL.Path
	vd.Yield["Standalone"] = buildinfo.Standalone

	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.ErrorWrap(err, fmt.Sprintf("executing template for URI '%s'", r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(v.App.HTTP500Page)
		return
	}

	w.WriteHeader(statusCode)
	io.Copy(w, &buf)
}
