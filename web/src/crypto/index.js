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

// module crypto.js provides cryptography operations using the Web Crypto API

import { utf8ToBuf, bufToB64, b64ToBuf } from '../libs/encoding';
import { PBKDF2, HKDF, SHA256, AES_GCM, AES_GCM_NONCE_SIZE } from './consts';

function mergeBuffers(buf1, buf2) {
  const buf = new ArrayBuffer(buf1.byteLength + buf2.byteLength);
  const bufView = new Uint8Array(buf);

  bufView.set(new Uint8Array(buf1), 0);
  bufView.set(new Uint8Array(buf2), buf1.byteLength);

  return buf;
}

// genRandomBytes returns an ArrayBuffer of random bytes of a given byte length.
function genRandomBytes(len) {
  const arr = new Uint8Array(len);
  return window.crypto.getRandomValues(arr);
}

export async function importAes256GcmKey(keyBuf) {
  const key = await window.crypto.subtle.importKey(
    'raw',
    keyBuf,
    { name: AES_GCM },
    false,
    ['encrypt', 'decrypt']
  );

  return key;
}

export async function aes256GcmEncrypt(keyBuf, dataBuf) {
  const cipherKey = await importAes256GcmKey(keyBuf);
  const iv = genRandomBytes(AES_GCM_NONCE_SIZE);

  const encrypted = await window.crypto.subtle.encrypt(
    {
      name: AES_GCM,
      iv
    },
    cipherKey,
    dataBuf
  );

  return mergeBuffers(iv, encrypted);
}

export async function aes256GcmDecrypt(keyBuf, dataBuf) {
  const cipherKey = await importAes256GcmKey(keyBuf);

  // split iv and ciphertext
  const ivBuf = dataBuf.slice(0, AES_GCM_NONCE_SIZE);
  const cipherTextBuf = dataBuf.slice(AES_GCM_NONCE_SIZE);

  const decrypted = await window.crypto.subtle.decrypt(
    {
      name: AES_GCM,
      iv: ivBuf
    },
    cipherKey,
    cipherTextBuf
  );

  return decrypted;
}

export async function pbkdf2(secretBuf, saltBuf, iterations) {
  const key = await window.crypto.subtle.importKey(
    'raw',
    secretBuf,
    { name: PBKDF2 },
    false,
    ['deriveBits']
  );

  return window.crypto.subtle.deriveBits(
    {
      name: PBKDF2,
      salt: saltBuf,
      iterations,
      hash: { name: SHA256 }
    },
    key,
    256
  );
}

export async function hkdf(secretBuf, saltBuf, infoBuf, algorithm, dkLen) {
  const key = await window.crypto.subtle.importKey(
    'raw',
    secretBuf,
    {
      name: HKDF
    },
    false,
    ['deriveBits']
  );

  return window.crypto.subtle.deriveBits(
    {
      name: HKDF,
      hash: algorithm,
      salt: saltBuf,
      info: infoBuf
    },
    key,
    dkLen
  );
}

// registerHelper generates and returns a set of keys for registration purposes
export async function registerHelper({ email, password, iteration }) {
  const emailBuf = utf8ToBuf(email);
  const passwordBuf = utf8ToBuf(password);

  const masterKey = await pbkdf2(passwordBuf, emailBuf, iteration);
  const authKey = await hkdf(
    masterKey,
    emailBuf,
    utf8ToBuf('auth'),
    SHA256,
    256
  );

  const cipherKey = genRandomBytes(32);
  const cipherKeyEnc = await aes256GcmEncrypt(masterKey, cipherKey);

  return {
    cipherKey: bufToB64(cipherKey),
    cipherKeyEnc: bufToB64(cipherKeyEnc),
    authKey: bufToB64(authKey)
  };
}

export async function loginHelper({ email, password, iteration }) {
  const emailBuf = utf8ToBuf(email);
  const passwordBuf = utf8ToBuf(password);

  const masterKey = await pbkdf2(passwordBuf, emailBuf, iteration);
  const authKey = await hkdf(
    masterKey,
    emailBuf,
    utf8ToBuf('auth'),
    SHA256,
    256
  );

  return {
    masterKey: bufToB64(masterKey),
    authKey: bufToB64(authKey)
  };
}

// getCipherKey returns cipher key in ArrayBuffer
export function getCipherKey(demo = false) {
  if (demo) {
    return b64ToBuf('demo');
  }

  const cipherKey = localStorage.getItem('cipherKey');
  const cipherKeyBuf = b64ToBuf(cipherKey);

  return cipherKeyBuf;
}
