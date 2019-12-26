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

import React, { useEffect } from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import Helmet from 'react-helmet';
import { Location } from 'history';

import { DigestNoteData } from 'jslib/operations/types';
import { parseSearchString } from 'jslib/helpers/url';
import { usePrevious } from 'web/libs/hooks';
import { Sort, SearchParams } from './types';
import { getDigest } from '../../store/digest';
import { useDispatch, useSelector } from '../../store';
import Header from './Header';
import Toolbar from './Toolbar';
import NoteList from './NoteList';
import Flash from '../Common/Flash';
import styles from './Digest.scss';

function useFetchData(digestUUID: string) {
  const dispatch = useDispatch();

  const { digest } = useSelector(state => {
    return {
      digest: state.digest
    };
  });

  const prevDigestUUID = usePrevious(digestUUID);

  useEffect(() => {
    if (!digest.isFetched || (digestUUID && prevDigestUUID !== digestUUID)) {
      dispatch(getDigest(digestUUID));
    }
  }, [dispatch, digestUUID, digest.isFetched, prevDigestUUID]);
}

interface Match {
  digestUUID: string;
}

interface Props extends RouteComponentProps<Match> {}

function getNotes(notes: DigestNoteData[], sort: string): DigestNoteData[] {
  return notes.concat().sort((i, j) => {
    if (sort === Sort.Oldest) {
      return new Date(i.createdAt).getTime() - new Date(j.createdAt).getTime();
    }

    return new Date(j.createdAt).getTime() - new Date(i.createdAt).getTime();
  });
}

function parseSearchParams(location: Location): SearchParams {
  const searchObj = parseSearchString(location.search);

  let sort: Sort;
  if (searchObj.sort === Sort.Oldest) {
    sort = Sort.Oldest;
  } else {
    sort = Sort.Newest;
  }

  return {
    sort,
    books: []
  };
}

const Digest: React.FunctionComponent<Props> = ({ location, match }) => {
  const { digestUUID } = match.params;

  useFetchData(digestUUID);

  const { digest } = useSelector(state => {
    return {
      digest: state.digest
    };
  });

  const params = parseSearchParams(location);
  const notes = getNotes(digest.data.notes, params.sort);

  return (
    <div className="page page-mobile-full">
      <Helmet>
        <title>Digest</title>
      </Helmet>

      <Header digest={digest.data} isFetched={digest.isFetched} />

      <div className="container mobile-nopadding">
        <Toolbar
          digestUUID={digest.data.uuid}
          sort={params.sort}
          isFetched={digest.isFetched}
        />

        <Flash
          kind="danger"
          when={digest.errorMessage !== null}
          wrapperClassName={styles['error-flash']}
        >
          Error getting digest: {digest.errorMessage}
        </Flash>
      </div>

      <div className="container mobile-nopadding">
        <NoteList
          digest={digest.data}
          notes={notes}
          isFetched={digest.isFetched}
          isFetching={digest.isFetching}
        />
      </div>
    </div>
  );
};

export default withRouter(Digest);
