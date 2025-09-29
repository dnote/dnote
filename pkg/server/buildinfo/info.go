/* Copyright (C) 2019, 2020, 2021, 2022, 2023, 2024, 2025 Dnote contributors
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package buildinfo

var (
	// Version is the server version
	Version = "master"
	// CSSFiles is the css files
	CSSFiles = ""
	// JSFiles is the js files
	JSFiles = ""
	// RootURL is the root url
	RootURL = "/"
	// Standalone reprsents whether the build is for on-premises. It is a string
	// rather than a boolean, so that it can be overridden during compile time.
	Standalone = "false"
)
