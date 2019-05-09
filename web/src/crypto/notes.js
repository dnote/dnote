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

import { aes256GcmDecrypt } from './index';
import { b64ToBuf, bufToUtf8 } from '../libs/encoding';

export async function decryptNote(note, cipherKeyBuf) {
  try {
    const contentDec = await aes256GcmDecrypt(
      cipherKeyBuf,
      b64ToBuf(note.content)
    );
    const bookLabelDec = await aes256GcmDecrypt(
      cipherKeyBuf,
      b64ToBuf(note.book.label)
    );

    return {
      ...note,
      content: bufToUtf8(contentDec),
      book: {
        ...note.book,
        label: bufToUtf8(bookLabelDec)
      }
    };
  } catch (e) {
    console.log(`Error while decrypting note ${note.uuid}`, e);
    console.log(e.stack);
    throw e;
  }
}
