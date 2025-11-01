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

// Package consts provides definitions of constants
package consts

var (
	// LegacyDnoteDirName is the name of the legacy directory containing dnote files
	LegacyDnoteDirName = ".dnote"
	// DnoteDirName is the name of the directory containing dnote files
	DnoteDirName = "dnote"
	// DnoteDBFileName is a filename for the Dnote SQLite database
	DnoteDBFileName = "dnote.db"
	// TmpContentFileBase is the base for the filename for a temporary content
	TmpContentFileBase = "DNOTE_TMPCONTENT"
	// TmpContentFileExt is the extension for the temporary content file
	TmpContentFileExt = "md"
	// ConfigFilename is the name of the config file
	ConfigFilename = "dnoterc"

	// SystemSchema is the key for schema in the system table
	SystemSchema = "schema"
	// SystemRemoteSchema is the key for remote schema in the system table
	SystemRemoteSchema = "remote_schema"
	// SystemLastSyncAt is the timestamp of the server at the last sync
	SystemLastSyncAt = "last_sync_time"
	// SystemLastMaxUSN is the user's max_usn from the server at the alst sync
	SystemLastMaxUSN = "last_max_usn"
	// SystemLastUpgrade is the timestamp at which the system more recently checked for an upgrade
	SystemLastUpgrade = "last_upgrade"
	// SystemSessionKey is the session key
	SystemSessionKey = "session_token"
	// SystemSessionKeyExpiry is the timestamp at which the session key will expire
	SystemSessionKeyExpiry = "session_token_expiry"
)
