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

package version

import (
	"fmt"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/spf13/cobra"
)

// NewCmd returns a new version command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Dnote",
		Long:  "Print the version number of Dnote",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dnote %s\n", ctx.Version)
		},
	}

	return cmd
}
