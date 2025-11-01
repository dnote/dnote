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
	"fmt"
	"html/template"
	"io/fs"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/pkg/errors"
)

// Engine is responsible for instantiating a new View
type Engine struct {
	filePatterns []string
	fileSystem   fs.FS
}

// NewEngine returns a new Engine
func NewEngine(filePatterns []string, fileSystem fs.FS) *Engine {
	return &Engine{
		filePatterns: filePatterns,
		fileSystem:   fileSystem,
	}
}

// NewDefaultEngine returns a new default Engine
func NewDefaultEngine() *Engine {
	patterns := []string{}

	patterns = append(patterns, iconFiles())
	patterns = append(patterns, layoutFiles())
	patterns = append(patterns, partialFiles())

	return NewEngine(patterns, TemplateFs)
}

// getTargetFiles returns an array of files needed for rendering
func (e Engine) getTargetFiles(files []string) []string {
	addTemplatePath(files)
	addTemplateExt(files)

	return append(files, e.filePatterns...)
}

// NewView returns a new view by parsing  the given layout and files
func (e Engine) NewView(app *app.App, viewConfig Config, files ...string) *View {
	viewHelpers := initHelpers(viewConfig, app)
	t := template.New(viewConfig.Title).Funcs(viewHelpers)

	targetFiles := e.getTargetFiles(files)

	t, err := t.ParseFS(e.fileSystem, targetFiles...)
	if err != nil {
		panic(errors.Wrap(err, "instantiating view"))
	}

	return &View{
		Template:    t,
		Layout:      viewConfig.getLayout(),
		AlertInBody: viewConfig.AlertInBody,
		App:         app,
	}
}

// layoutFiles returns a slice of strings representing
// the layout files used in our application.
func layoutFiles() string {
	return fmt.Sprintf("templates/layouts/*%s", TemplateExt)
}

// iconFiles returns a slice of strings representing
// the icon files used in our application.
func iconFiles() string {
	return fmt.Sprintf("templates/icons/*%s", TemplateExt)
}

func partialFiles() string {
	return fmt.Sprintf("templates/partials/*%s", TemplateExt)
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates.
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = fmt.Sprintf("templates/%s", f)
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
		files[i] = f + TemplateExt
	}
}
