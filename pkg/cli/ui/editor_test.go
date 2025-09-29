/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
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
		ctx := context.InitTestCtx(t, context.Paths{
			Data:  "../tmp",
			Cache: "../tmp",
		}, nil)
		defer context.TeardownTestCtx(t, ctx)

		res, err := GetTmpContentPath(ctx)
		if err != nil {
			t.Fatal(errors.Wrap(err, "executing"))
		}

		expected := fmt.Sprintf("%s/%s", ctx.Paths.Cache, "DNOTE_TMPCONTENT_0.md")
		assert.Equal(t, res, expected, "filename did not match")
	})

	t.Run("one existing session", func(t *testing.T) {
		// set up
		ctx := context.InitTestCtx(t, context.Paths{
			Data:  "../tmp2",
			Cache: "../tmp2",
		}, nil)
		defer context.TeardownTestCtx(t, ctx)

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
		ctx := context.InitTestCtx(t, context.Paths{
			Data:  "../tmp3",
			Cache: "../tmp3",
		}, nil)
		defer context.TeardownTestCtx(t, ctx)

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
