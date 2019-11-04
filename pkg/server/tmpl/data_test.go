package tmpl

import (
	"html/template"
	"testing"
	"time"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/server/database"
	"github.com/pkg/errors"
)

func TestDefaultPageGetData(t *testing.T) {
	p := defaultPage{}

	result := p.getData()

	assert.Equal(t, result.MetaTags, template.HTML(""), "MetaTags mismatch")
	assert.Equal(t, result.Title, "Dnote", "Title mismatch")
}

func TestNotePageGetData(t *testing.T) {
	a, err := NewAppShell(nil)
	if err != nil {
		t.Fatal(errors.Wrap(err, "preparing app shell"))
	}

	p := notePage{
		Note: database.Note{
			Book: database.Book{
				Label: "vocabulary",
			},
			AddedOn: time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC).UnixNano(),
		},
		T: a.T,
	}

	result, err := p.getData()
	if err != nil {
		t.Fatal(errors.Wrap(err, "executing"))
	}

	assert.NotEqual(t, result.MetaTags, template.HTML(""), "MetaTags should not be empty")
	assert.Equal(t, result.Title, "Note: vocabulary (Jan 2 2019)", "Title mismatch")
}
