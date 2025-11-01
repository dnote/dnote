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

var apiEndpointFlag string

// NewCmd returns a new logout command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "Logout from the server",
		Example: example,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVar(&apiEndpointFlag, "apiEndpoint", "", "API endpoint to connect to (defaults to value in config)")

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
		// Override APIEndpoint if flag was provided
		if apiEndpointFlag != "" {
			ctx.APIEndpoint = apiEndpointFlag
		}

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
