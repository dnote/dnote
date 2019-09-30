import { Location } from 'history';

import { parseSearchString } from 'jslib/helpers/url';
import { removeKey } from 'jslib/helpers/obj';
import { Queries, keywordBook } from 'jslib/helpers/queries';
import { getHomePath } from './paths';

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
