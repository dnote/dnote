/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

// module ext provides a cross-browser interface to access extension APIs
// by using WebExtensions API if available, and using Chrome as a fallback.
const ext: any = {};

const apis = ['tabs', 'storage', 'runtime'];

for (let i = 0; i < apis.length; i++) {
  const api = apis[i];

  try {
    if (browser[api]) {
      ext[api] = browser[api];
    }
  } catch (e) {}

  try {
    if (chrome[api] && !ext[api]) {
      ext[api] = chrome[api];

      // Standardize the signature to conform to WebExtensions API
      if (api === 'tabs') {
        const fn = ext[api].create;

        // Promisify chrome.tabs.create
        ext[api].create = obj => {
          return new Promise(resolve => {
            fn(obj, tab => {
              resolve(tab);
            });
          });
        };
      }
    }
  } catch (e) {}
}

export default ext;
