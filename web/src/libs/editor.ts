/* Copyright (C) 2019 Monomax Software Pty Ltd
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

// sessionKeyNew is the editor session key for a new note
const sessionKeyNew = 'new';

// getEditorSessionkey returns a unique editor session key for the given noteUUID.
// If the noteUUID is null, it returns a session key for the new note.
// Editor session holds an editor state for a particular note.
export function getEditorSessionkey(noteUUID: string | null): string {
  if (noteUUID === null) {
    return sessionKeyNew;
  }

  return noteUUID;
}
