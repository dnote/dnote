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

package root

import (
	"github.com/spf13/cobra"
)

var apiEndpointFlag string

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
	root.PersistentFlags().StringVar(&apiEndpointFlag, "api-endpoint", "", "the API endpoint to connect to (defaults to value in config)")
}

// GetRoot returns the root command
func GetRoot() *cobra.Command {
	return root
}

// GetAPIEndpointFlag returns the value of the --api-endpoint flag
func GetAPIEndpointFlag() string {
	return apiEndpointFlag
}

// Register adds a new command
func Register(cmd *cobra.Command) {
	root.AddCommand(cmd)
}

// Execute runs the main command
func Execute() error {
	return root.Execute()
}
