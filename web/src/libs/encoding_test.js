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
import { utf8ToBuf, bufToUtf8, bufToB64, b64ToBuf } from './encoding';

describe('utf8ToBuf', () => {
  it('converts a string to an ArrayBuffer', () => {
    const result = utf8ToBuf('hi');
    const buf = new ArrayBuffer(2);
    const bufView = new Uint8Array(buf);

    bufView[0] = 104;
    bufView[1] = 105;
    expect(result).to.deep.equal(buf);
  });
});

describe('bufToUtf8', () => {
  it('converts an ArrayBuffer to a UTF-8 string', () => {
    const str = 'hi';
    const buf = new ArrayBuffer(2);
    const bufView = new Uint8Array(buf);
    bufView[0] = str.charCodeAt(0);
    bufView[1] = str.charCodeAt(1);

    const result = bufToUtf8(buf);
    expect(result).to.equal(str);
  });
});

describe('bufToB64', () => {
  it('encodes the given ArrayBuffer using base64', () => {
    const str = 'hi';

    const buf = new ArrayBuffer(2);
    const bufView = new Uint8Array(buf);
    bufView[0] = str.charCodeAt(0);
    bufView[1] = str.charCodeAt(1);

    const result = bufToB64(buf);
    expect(result).to.equal('aGk=');
  });
});

describe('b64ToBuf', () => {
  it('converts the given base64 string into an ArrayBuffer', () => {
    const buf = new ArrayBuffer(2);
    const bufView = new Uint8Array(buf);
    bufView[0] = 104;
    bufView[1] = 105;

    const str = 'aGk=';
    const result = b64ToBuf(str);
    expect(result).to.deep.equal(buf);
  });
});
