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

// This module provides interfaces to perform operations. It abstarcts
// the backend implementation and thus unifies the API for web and desktop clients.

import * as notesService from '../services/notes';
import { aes256GcmEncrypt, getCipherKey } from '../crypto';
import { decryptNote } from '../crypto/notes';
import { utf8ToBuf, bufToB64 } from '../libs/encoding';

// create creates an encrypted note. It returns a promise that resolves with
// a decrypted note
export async function create({ bookUUID, content }) {
  const cipherKeyBuf = getCipherKey();

  const contentBuf = utf8ToBuf(content);
  const contentEnc = await aes256GcmEncrypt(cipherKeyBuf, contentBuf);

  return notesService
    .create({ bookUUID, content: bufToB64(contentEnc) })
    .then(response => {
      const note = response.result;

      return decryptNote(note, cipherKeyBuf);
    });
}

export async function update(noteUUID, input) {
  const cipherKeyBuf = getCipherKey();

  const payload = {
    ...input
  };

  if (input.content !== undefined) {
    const contentBuf = utf8ToBuf(input.content);
    const contentEnc = await aes256GcmEncrypt(cipherKeyBuf, contentBuf);
    payload.content = bufToB64(contentEnc);
  }

  return notesService.update(noteUUID, payload).then(response => {
    const note = response.result;

    return decryptNote(note, cipherKeyBuf);
  });
}

export async function remove(noteUUID) {
  return notesService.remove(noteUUID);
}
