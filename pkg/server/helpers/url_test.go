package helpers

import (
	"net/url"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
)

func TestGetPath(t *testing.T) {
	t.Run("without query", func(t *testing.T) {
		// execute
		got := GetPath("/some-path", nil)

		// test
		assert.Equal(t, got, "/some-path", "got mismatch")
	})

	t.Run("with query", func(t *testing.T) {
		// execute
		q := url.Values{}
		q.Set("foo", "bar")
		q.Set("baz", "/quz")
		got := GetPath("/some-path", &q)

		// test
		assert.Equal(t, got, "/some-path?baz=%2Fquz&foo=bar", "got mismatch")
	})
}
