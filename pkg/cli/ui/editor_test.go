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

package ui

import (
	"fmt"
	"os"
	"testing"

	"github.com/dnote/dnote/pkg/assert"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/pkg/errors"
)

func TestGetTmpContentPath(t *testing.T) {
	t.Run("no collision", func(t *testing.T) {
		ctx := context.InitTestCtx(t)

		res, err := GetTmpContentPath(ctx)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		expected := fmt.Sprintf("%s/%s", ctx.Paths.Cache, "DNOTE_TMPCONTENT_0.md")
		assert.Equal(t, res, expected, "filename did not match")
	})

	t.Run("one existing session", func(t *testing.T) {
		// set up
		ctx := context.InitTestCtx(t)

		p := fmt.Sprintf("%s/%s", ctx.Paths.Cache, "DNOTE_TMPCONTENT_0.md")
		if _, err := os.Create(p); err != nil {
			t.Fatal(errors.Wrap(err, "preparing the conflicting file"))
		}

		// execute
		res, err := GetTmpContentPath(ctx)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		// test
		expected := fmt.Sprintf("%s/%s", ctx.Paths.Cache, "DNOTE_TMPCONTENT_1.md")
		assert.Equal(t, res, expected, "filename did not match")
	})

	t.Run("two existing sessions", func(t *testing.T) {
		// set up
		ctx := context.InitTestCtx(t)

		p1 := fmt.Sprintf("%s/%s", ctx.Paths.Cache, "DNOTE_TMPCONTENT_0.md")
		if _, err := os.Create(p1); err != nil {
			t.Fatal(errors.Wrap(err, "preparing the conflicting file"))
		}
		p2 := fmt.Sprintf("%s/%s", ctx.Paths.Cache, "DNOTE_TMPCONTENT_1.md")
		if _, err := os.Create(p2); err != nil {
			t.Fatal(errors.Wrap(err, "preparing the conflicting file"))
		}

		// execute
		res, err := GetTmpContentPath(ctx)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		// test
		expected := fmt.Sprintf("%s/%s", ctx.Paths.Cache, "DNOTE_TMPCONTENT_2.md")
		assert.Equal(t, res, expected, "filename did not match")
	})
}
