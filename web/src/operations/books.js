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

import * as booksService from '../services/books';
import { aes256GcmEncrypt, getCipherKey } from '../crypto';
import { decryptBook } from '../crypto/books';
import { utf8ToBuf, bufToB64 } from '../libs/encoding';

export async function get(bookUUID) {
  const cipherKeyBuf = getCipherKey();

  return booksService.get(bookUUID).then(book => {
    return decryptBook(book, cipherKeyBuf);
  });
}

// create creates an encrypted book. It returns a promise that resolves with
// a decrypted book.
export async function create({ name }) {
  const cipherKeyBuf = getCipherKey();

  const nameBuf = utf8ToBuf(name);
  const nameEnc = await aes256GcmEncrypt(cipherKeyBuf, nameBuf);

  return booksService.create({ name: bufToB64(nameEnc) }).then(response => {
    const { book } = response;

    return decryptBook(book, cipherKeyBuf);
  });
}

// remove deletes the book with the given uuid
export async function remove(bookUUID) {
  return booksService.remove(bookUUID);
}
