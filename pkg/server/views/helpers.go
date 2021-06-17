package views

import (
	"strconv"
	"strings"
	"time"

	"github.com/dnote/dnote/pkg/server/buildinfo"
)

func toDateTime(year, month int) string {
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

func getFullMonthName(month int) string {
	return time.Month(month).String()
}

func rootURL() string {
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
func excerpt(s string, maxLength int) string {
	if len(s) < maxLength {
		return s
	}

	ret := s[0 : maxLength+1]

	last := max(0, min(len(ret), strings.LastIndex(ret, " ")))

	ret = ret[0:last]
	ret += "..."

	return ret
}
