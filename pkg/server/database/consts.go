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

package database

const (
	// TokenTypeResetPassword is a type of a token for reseting password
	TokenTypeResetPassword = "reset_password"
)

const (
	// BookDomainAll incidates that all books are eligible to be the source books
	BookDomainAll = "all"
	// BookDomainIncluding incidates that some specified books are eligible to be the source books
	BookDomainIncluding = "including"
	// BookDomainExluding incidates that all books except for some specified books are eligible to be the source books
	BookDomainExluding = "excluding"
)
