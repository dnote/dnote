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

type FilterFn = (val: any, key: any) => boolean;

export function filterObjKeys(obj: object, filterFn: FilterFn) {
  return Object.keys(obj)
    .filter(key => {
      const val = obj[key];
      return filterFn(val, key);
    })
    .reduce((ret, key) => {
      return {
        ...ret,
        [key]: obj[key]
      };
    }, {});
}

// whitelist returns a new object whose keys are whitelisted by the given array
// of keys
export function whitelist(obj, keys) {
  return filterObjKeys(obj, (val, key) => {
    return keys.indexOf(key) > -1;
  });
}

// blacklist returns a new object where key-val pairs are filtered out by keys
export function blacklist(obj, keys) {
  return filterObjKeys(obj, (val, key) => {
    return keys.indexOf(key) === -1;
  });
}

// isEmptyObj checks that an object does not have any properties of its own
export function isEmptyObj(obj) {
  return Object.getOwnPropertyNames(obj).length === 0;
}

// removeKey returns a new object with the given key removed
export function removeKey(obj: object, deleteKey: string) {
  const keys = Object.keys(obj).filter(key => key !== deleteKey);

  const ret = {};

  for (let i = 0; i < keys.length; ++i) {
    const key = keys[i];
    ret[key] = obj[key];
  }

  return ret;
}
