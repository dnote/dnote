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

package root

import (
	"github.com/spf13/cobra"
)

var dbPathFlag string

var root = &cobra.Command{
	Use:           "dnote",
	Short:         "Dnote - a simple command line notebook",
	SilenceErrors: true,
	SilenceUsage:  true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	root.PersistentFlags().StringVar(&dbPathFlag, "dbPath", "", "the path to the database file (defaults to standard location)")
}

// GetRoot returns the root command
func GetRoot() *cobra.Command {
	return root
}

// GetDBPathFlag returns the value of the --dbPath flag
func GetDBPathFlag() string {
	return dbPathFlag
}

// Register adds a new command
func Register(cmd *cobra.Command) {
	root.AddCommand(cmd)
}

// Execute runs the main command
func Execute() error {
	return root.Execute()
}
