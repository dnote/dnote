/* Copyright (C) 2019 Monomax Software Pty Ltd
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

package logout

import (
	"database/sql"

	"github.com/dnote/dnote/pkg/cli/client"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ErrNotLoggedIn is an error for logging out when not logged in
var ErrNotLoggedIn = errors.New("not logged in")

var example = `
  dnote logout`

// NewCmd returns a new logout command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "Logout from the server",
		Example: example,
		RunE:    newRun(ctx),
	}

	return cmd
}

// Do performs logout
func Do(ctx context.DnoteCtx) error {
	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	var key string
	err = database.GetSystem(tx, consts.SystemSessionKey, &key)
	if errors.Cause(err) == sql.ErrNoRows {
		return ErrNotLoggedIn
	} else if err != nil {
		return errors.Wrap(err, "getting session key")
	}

	err = client.Signout(ctx, key)
	if err != nil {
		return errors.Wrap(err, "requesting logout")
	}

	if err := database.DeleteSystem(tx, consts.SystemSessionKey); err != nil {
		return errors.Wrap(err, "deleting session key")
	}
	if err := database.DeleteSystem(tx, consts.SystemSessionKeyExpiry); err != nil {
		return errors.Wrap(err, "deleting session key expiry")
	}

	tx.Commit()

	return nil
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		err := Do(ctx)
		if err == ErrNotLoggedIn {
			log.Error("not logged in\n")
			return nil
		} else if err != nil {
			return errors.Wrap(err, "logging out")
		}

		log.Success("logged out\n")

		return nil
	}
}
