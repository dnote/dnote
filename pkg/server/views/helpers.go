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
	"strconv"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/app"
	"github.com/dnote/dnote/pkg/server/buildinfo"
	"github.com/pkg/errors"
	"html/template"
)

func initHelpers(c Config, a *app.App) template.FuncMap {
	ctx := newViewCtx(c)

	ret := template.FuncMap{
		"csrfField":        ctx.csrfField,
		"css":              ctx.css,
		"js":               ctx.js,
		"title":            ctx.title,
		"headerTemplate":   ctx.headerTemplate,
		"rootURL":          ctx.rootURL,
		"getFullMonthName": ctx.getFullMonthName,
		"toDateTime":       ctx.toDateTime,
		"excerpt":          ctx.excerpt,
		"timeAgo":          ctx.timeAgo,
		"timeFormat":       ctx.timeFormat,
		"toISOString":      ctx.toISOString,
		"dict":             ctx.dict,
		"defaultValue":     ctx.defaultValue,
		"add":              ctx.add,
		"assetBaseURL": func() string {
			return a.AssetBaseURL
		},
	}

	// extend with helpers that are defined specific to a view
	if c.HelperFuncs != nil {
		for k, v := range c.HelperFuncs {
			ret[k] = v
		}
	}

	return ret
}

func (v viewCtx) csrfField() (template.HTML, error) {
	return "", errors.New("csrfField is not implemented")
}

func (v viewCtx) css() []string {
	return strings.Split(buildinfo.CSSFiles, ",")
}

func (v viewCtx) js() []string {
	return strings.Split(buildinfo.JSFiles, ",")
}

func (v viewCtx) title() string {
	if v.Config.Title != "" {
		return fmt.Sprintf("%s | %s", v.Config.Title, siteTitle)
	}

	return siteTitle
}

func (v viewCtx) headerTemplate() string {
	return v.Config.HeaderTemplate
}

func (v viewCtx) toDateTime(year, month int) string {
	sb := strings.Builder{}

	sb.WriteString(strconv.Itoa(year))
	sb.WriteString("-")

	if month < 10 {
		sb.WriteString("0")
		sb.WriteString(strconv.Itoa(month))
	} else {
		sb.WriteString(strconv.Itoa(month))
	}

	return sb.String()
}

func (v viewCtx) getFullMonthName(month int) string {
	return time.Month(month).String()
}

func (v viewCtx) rootURL() string {
	return buildinfo.RootURL
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// excerpt trims the given string up to the last word that makes the string
// exceed the maxLength, and attaches ellipses at the end. If the string is
// shorter than the given maxLength, it returns the original string.
func (v viewCtx) excerpt(s string, maxLength int) string {
	if len(s) < maxLength {
		return s
	}

	ret := s[0:maxLength]
	ret = s[0:min(len(ret), max(0, strings.LastIndex(ret, " ")))]
	ret += "..."

	return ret
}

func (v viewCtx) timeFormat(t time.Time, format string) string {
	return t.Format(format)
}

func (v viewCtx) timeAgo(t time.Time) string {
	now := v.Clock.Now()
	diff := relativeTimeDiff(now, t)

	if diff.tense == "past" {
		return fmt.Sprintf("%s ago", diff.text)
	}

	if diff.tense == "future" {
		return fmt.Sprintf("in %s", diff.text)
	}

	return diff.text
}

func (v viewCtx) toISOString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func (v viewCtx) dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func (v viewCtx) defaultValue(value, fallback interface{}) interface{} {
	if value == nil {
		return fallback
	}

	return value
}

func (v viewCtx) add(a, b int) interface{} {
	return a + b
}
