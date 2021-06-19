package views

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/buildinfo"
	"github.com/pkg/errors"
	"html/template"
)

func initHelpers(c Config) template.FuncMap {
	ctx := newViewCtx(c)

	ret := template.FuncMap{
		"csrfField":        ctx.csrfField,
		"css":              ctx.css,
		"title":            ctx.title,
		"headerTemplate":   ctx.headerTemplate,
		"rootURL":          ctx.rootURL,
		"getFullMonthName": ctx.getFullMonthName,
		"toDateTime":       ctx.toDateTime,
		"excerpt":          ctx.excerpt,
		"timeAgo":          ctx.timeAgo,
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

	ret := s[0 : maxLength+1]

	last := max(0, min(len(ret), strings.LastIndex(ret, " ")))

	ret = ret[0:last]
	ret += "..."

	return ret
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
