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

// import { receiveNote } from '../note';
// import { refreshNoteInList } from '../notes';

import {
  FLUSH_CONTENT,
  UPDATE_BOOK,
  RESET,
  STAGE_NOTE,
  MARK_DIRTY,
  MarkDirtyAction,
  StageNoteAction,
  FlushContentAction,
  UpdateBookAction,
  ResetAction
} from './type';

export function stageNote({
  noteUUID,
  bookUUID,
  bookLabel,
  content
}): StageNoteAction {
  return {
    type: STAGE_NOTE,
    data: { noteUUID, bookUUID, bookLabel, content }
  };
}

export function flushContent(content): FlushContentAction {
  return {
    type: FLUSH_CONTENT,
    data: { content }
  };
}

export interface UpdateBookActionParam {
  uuid: string;
  label: string;
}

export function updateBook({
  uuid,
  label
}: UpdateBookActionParam): UpdateBookAction {
  return {
    type: UPDATE_BOOK,
    data: {
      uuid,
      label
    }
  };
}

export function resetEditor(): ResetAction {
  return {
    type: RESET
  };
}

export function markDirty(): MarkDirtyAction {
  return {
    type: MARK_DIRTY
  };
}
