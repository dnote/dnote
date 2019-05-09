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

import * as notesOperation from '../operations/notes';
import { receiveNote } from './note';
import { refreshNoteInList } from './notes';

export const MARK_DIRTY = 'editor/MARK_DIRTY';
export const STAGE_NOTE = 'editor/STAGE_NOTE';
export const COMMIT_NOTE = 'editor/COMMIT_NOTE';
export const UPDATE_CONTENT = 'editor/UPDATE_CONTENT';
export const UPDATE_BOOK_UUID = 'editor/UPDATE_BOOK_UUID';
export const RESET = 'editor/RESET';

export function stageNote({ noteUUID, bookUUID, content }) {
  return {
    type: STAGE_NOTE,
    data: { noteUUID, bookUUID, content }
  };
}

export function commitNote() {
  return (dispatch, getState) => {
    const { editor } = getState();

    return notesOperation
      .update(editor.noteUUID, {
        content: editor.content,
        bookUUID: editor.bookUUID
      })
      .then(note => {
        dispatch(receiveNote(note));
        dispatch(
          refreshNoteInList({ noteUUID: note.uuid, refreshedNote: note })
        );
        dispatch({
          type: COMMIT_NOTE
        });
      })
      .catch(err => {
        // TODO: if deleted elsewhere, create
        console.log('err', err);
      });
  };
}

export function updateContent(content) {
  return {
    type: UPDATE_CONTENT,
    data: { content }
  };
}

export function updateBookUUID(uuid) {
  return {
    type: UPDATE_BOOK_UUID,
    data: {
      uuid
    }
  };
}

export function resetEditor() {
  return {
    type: RESET
  };
}

export function markDirty() {
  return {
    type: MARK_DIRTY
  };
}
