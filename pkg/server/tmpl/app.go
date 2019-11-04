package tmpl

import (
	"bytes"
	"html/template"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
)

// routes
var notesPathRegex = regexp.MustCompile("^/notes/([^/]+)$")

// template names
var templateIndex = "index"
var templateNoteMetaTags = "note_metatags"

// AppShell represents the application in HTML
type AppShell struct {
	T *template.Template
}

// ErrNotFound is an error indicating that a resource was not found
var ErrNotFound = errors.New("not found")

// NewAppShell parses the templates for the application
func NewAppShell(content []byte) (AppShell, error) {
	t, err := template.New(templateIndex).Parse(string(content))
	if err != nil {
		return AppShell{}, errors.Wrap(err, "parsing the index template")
	}

	_, err = t.New(templateNoteMetaTags).Parse(noteMetaTags)
	if err != nil {
		return AppShell{}, errors.Wrap(err, "parsing the note meta tags template")
	}

	return AppShell{t}, nil
}

// Execute executes the index template
func (a AppShell) Execute(r *http.Request) ([]byte, error) {
	data, err := a.getData(r)

	var buf bytes.Buffer
	if err != nil {
		return nil, errors.Wrap(err, "getting data")
	}

	if err := a.T.ExecuteTemplate(&buf, templateIndex, data); err != nil {
		return nil, errors.Wrap(err, "executing template")
	}

	return buf.Bytes(), nil
}

func (a AppShell) getData(r *http.Request) (tmplData, error) {
	path := r.URL.Path

	if ok, params := matchPath(path, notesPathRegex); ok {
		p, err := a.newNotePage(r, params[0])
		if err != nil {
			return tmplData{}, errors.Wrap(err, "instantiating note page")
		}

		return p.getData()
	}

	p := defaultPage{}
	return p.getData(), nil
}

// matchPath checks if the given path matches the given regular expressions
// and returns a boolean as well as any parameters from regex capture groups.
func matchPath(p string, reg *regexp.Regexp) (bool, []string) {
	match := notesPathRegex.FindStringSubmatch(p)

	if len(match) > 0 {
		return true, match[1:]
	}

	return false, nil
}
