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

export const AES_GCM = 'AES-GCM';
export const PBKDF2 = 'PBKDF2';
export const HKDF = 'HKDF';
export const SHA256 = 'SHA-256';
export const DEFAULT_KDF_ITERATION = 100000;

// AES_GCM_NONCE_SIZE is the size of the iv, in bytes, of AES in GCM mode
export const AES_GCM_NONCE_SIZE = 12;
