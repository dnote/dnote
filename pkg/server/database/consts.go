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
