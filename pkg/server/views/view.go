package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/config"
	"github.com/dnote/dnote/pkg/server/context"
	"github.com/dnote/dnote/pkg/server/log"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
)

const (
	// templateExt is the template extension
	templateExt string = ".gohtml"
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

// NewView returns a new view by parsing  the given layout and files
func NewView(appConfig config.Config, viewConfig Config, files ...string) *View {
	baseDir := appConfig.PageTemplateDir
	addTemplatePath(baseDir, files)
	addTemplateExt(files)

	files = append(files, iconFiles(baseDir)...)
	files = append(files, layoutFiles(baseDir)...)
	files = append(files, partialFiles(baseDir)...)

	viewHelpers := initHelpers(viewConfig)
	t := template.New(viewConfig.Title).Funcs(viewHelpers)

	t, err := t.ParseFiles(files...)
	if err != nil {
		panic(errors.Wrap(err, "instantiating view."))
	}

	return &View{
		Template:    t,
		Layout:      viewConfig.getLayout(),
		AlertInBody: viewConfig.AlertInBody,
		StaticDir:   appConfig.StaticDir,
	}
}

// View holds the information about a view
type View struct {
	Template *template.Template
	Layout   string
	// AlertInBody specifies if alert should be set in the body instead of the header
	AlertInBody bool
	StaticDir   string
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

	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.ErrorWrap(err, fmt.Sprintf("executing template: '%s' at '%s'", v.Template.Name(), r.RequestURI))
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, fmt.Sprintf("%s/500.html", v.StaticDir))
		return
	}

	w.WriteHeader(statusCode)
	io.Copy(w, &buf)
}

func getFiles(pattern string) []string {
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	return files
}

// layoutFiles returns a slice of strings representing
// the layout files used in our application.
func layoutFiles(baseDir string) []string {
	return getFiles(fmt.Sprintf("%s/layouts/*%s", baseDir, templateExt))
}

// iconFiles returns a slice of strings representing
// the icon files used in our application.
func iconFiles(baseDir string) []string {
	return getFiles(fmt.Sprintf("%s/icons/*%s", baseDir, templateExt))
}

func partialFiles(baseDir string) []string {
	return getFiles(fmt.Sprintf("%s/partials/*%s", baseDir, templateExt))
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates.
func addTemplatePath(baseDir string, files []string) {
	for i, f := range files {
		files[i] = fmt.Sprintf("%s/%s", baseDir, f)
	}
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates and it appends
// the templateExt extension to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"home.gohtml"} if templateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + templateExt
	}
}
