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

// utf8ToBuf turns a given string into an ArrayBuffer
export function utf8ToBuf(str) {
  // convert to utf8 encoding
  const strUtf8 = unescape(encodeURIComponent(str));

  const buf = new ArrayBuffer(strUtf8.length);
  const bufView = new Uint8Array(buf);
  for (let i = 0; i < strUtf8.length; i++) {
    bufView[i] = strUtf8.charCodeAt(i);
  }

  return buf;
}

// bufToUtf8 turns a given ArrayBuffer into a UTF-8 encoded string
export function bufToUtf8(buf) {
  const bytes = new Uint8Array(buf);
  const encodedString = String.fromCharCode.apply(null, bytes);

  return decodeURIComponent(escape(encodedString));
}

// bufToB64 encodes the given ArrayBuffer using base64
export function bufToB64(buf) {
  let binary = '';
  const bytes = new Uint8Array(buf);
  const len = bytes.byteLength;

  for (let i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i]);
  }

  return window.btoa(binary);
}

// b64ToBuf turns a given base64 tring into an ArrayBuffer
export function b64ToBuf(base64Str) {
  const binary = window.atob(base64Str);
  const len = binary.length;

  const buf = new ArrayBuffer(len);
  const bytes = new Uint8Array(buf);
  for (let i = 0; i < len; i++) {
    bytes[i] = binary.charCodeAt(i);
  }

  return buf;
}
