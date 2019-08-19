import { Location } from 'history';

import { parseSearchString } from './url';
import { removeKey } from './obj';
import { getHomePath } from './paths';
import * as searchLib from './search';

export interface Queries {
  q: string;
  book: string[];
}

function encodeQuery(keyword: string, value: string): string {
  return `${keyword}:${value}`;
}

const keywordBook = 'book';

export const keywords = [keywordBook];

// parse unmarshals the given string represesntation of the queries into an object
export function parse(s: string): Queries {
  const result = searchLib.parse(s, keywords);

  let bookValue: string[];
  const { book } = result.filters;
  if (!book) {
    bookValue = [];
  } else if (typeof book === 'string') {
    bookValue = [book];
  } else {
    bookValue = book;
  }

  let qValue: string;
  if (result.text) {
    qValue = result.text;
  } else {
    qValue = '';
  }

  return {
    q: qValue,
    book: bookValue
  };
}

// stringify marshals the givne queries into a string format
export function stringify(queries: Queries): string {
  let ret = '';

  if (queries.book.length > 0) {
    for (let i = 0; i < queries.book.length; i++) {
      const book = queries.book[i];

      ret += encodeQuery(keywordBook, book);
      ret += ' ';
    }
  }

  if (queries.q !== '') {
    const result = searchLib.parse(queries.q, keywords);
    ret += result.text;
  }

  return ret;
}

export function getSearchDest(location: Location, queries: Queries) {
  let searchObj: any = parseSearchString(location.search);

  if (queries.q !== '') {
    searchObj.q = queries.q;
  } else {
    searchObj = removeKey(searchObj, 'q');
  }

  if (queries.book.length > 0) {
    searchObj.book = queries.book;
  } else {
    searchObj = removeKey(searchObj, keywordBook);
  }

  searchObj = removeKey(searchObj, 'page');

  return getHomePath(searchObj);
}
