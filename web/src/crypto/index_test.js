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

import { expect } from 'chai';
import {
  pbkdf2,
  hkdf,
  importAes256GcmKey,
  aes256GcmEncrypt,
  aes256GcmDecrypt
} from './index';
import { bufToB64, b64ToBuf, utf8ToBuf, bufToUtf8 } from '../libs/encoding';

describe('pbkdf2', () => {
  it('derives correct bits using SHA256', async () => {
    const secret = utf8ToBuf('v2biqXbuXabsuZWXXyQ76f7SvhxJxzpp');
    const salt = utf8ToBuf('tkcZv7RDyebPnD9DLt63kAxtIvHxcTe6');
    const key = await pbkdf2(secret, salt, 1000);

    expect(bufToB64(key)).to.equal(
      '8mR5rKfEAg+9LuF9SttGGV5yiOHCiQT/PQ1HORdVXYU='
    );
  });
});

describe('hkdf', () => {
  it('derives correct 256 bit key using SHA256', async () => {
    const secret = utf8ToBuf('v2biqXbuXabsuZWXXyQ76f7SvhxJxzpp');
    const salt = utf8ToBuf('tkcZv7RDyebPnD9DLt63kAxtIvHxcTe6');
    const info = utf8ToBuf('zwm3r2mmZ665m77cPQRU2hmPrzhrgQjg');
    const dkLen = 256;

    const key = await hkdf(secret, salt, info, 'SHA-256', dkLen);
    expect(bufToB64(key)).to.equal(
      'l7bgmvf9YlaDt+QJWca3XATnR1GCoh6XvDh/1/bVs5Y='
    );
  });

  it('derives correct 256 bit key using SHA256, given an ArrayBuffer as a secret', async () => {
    const secret = utf8ToBuf('v2biqXbuXabsuZWXXyQ76f7SvhxJxzpp');
    const salt = utf8ToBuf('tkcZv7RDyebPnD9DLt63kAxtIvHxcTe6');
    const info = utf8ToBuf('zwm3r2mmZ665m77cPQRU2hmPrzhrgQjg');
    const dkLen = 256;

    const key = await hkdf(secret, salt, info, 'SHA-256', dkLen);
    expect(bufToB64(key)).to.equal(
      'l7bgmvf9YlaDt+QJWca3XATnR1GCoh6XvDh/1/bVs5Y='
    );
  });
});

describe('importAes256GcmKey', () => {
  it('imports a key', async () => {
    // execute
    const secret = b64ToBuf('79fkmXp1Eu+O+1IBqHjDvcwciJM4k+nO9bEjj8bvSXo=');
    const key = await importAes256GcmKey(secret);

    // test
    expect(key).to.be.a('CryptoKey');
  });
});

describe('aes256Gcm', () => {
  const testCases = [
    { input: 'hello world' },
    { input: 'foo\nbar\nbaz 123\nquz' },
    { input: 'föo\nbār\nbåz & qūz' }
  ];

  for (let i = 0; i < testCases.length; i++) {
    const tc = testCases[i];

    it(`can encrypt and decrypt - test case ${i}`, async () => {
      const keyBuf = b64ToBuf('79fkmXp1Eu+O+1IBqHjDvcwciJM4k+nO9bEjj8bvSXo=');

      try {
        const inputBuf = utf8ToBuf(tc.input);
        const ciphertext = await aes256GcmEncrypt(keyBuf, inputBuf);
        const decoded = await aes256GcmDecrypt(keyBuf, ciphertext);

        expect(bufToUtf8(decoded)).to.equal(tc.input);
      } catch (e) {
        console.log(e.message);
        throw e;
      }
    });
  }
});
