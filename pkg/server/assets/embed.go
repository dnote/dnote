package assets

import (
	"embed"
	"github.com/pkg/errors"
	"io/fs"
)

//go:embed static
var staticFiles embed.FS

// GetStaticFS returns a filesystem for static files, with
// all files situated in the root of the filesystem
func GetStaticFS() (fs.FS, error) {
	subFs, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, errors.Wrap(err, "getting sub filesystem")
	}

	return subFs, nil
}

// MustGetHTTP500ErrorPage returns the content of HTML file for HTTP 500 error
func MustGetHTTP500ErrorPage() []byte {
	ret, err := staticFiles.ReadFile("static/500.html")
	if err != nil {
		panic(errors.Wrap(err, "reading HTML file for 500 HTTP error"))
	}

	return ret
}
