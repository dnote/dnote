/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import { usePrevious } from 'web/libs/hooks';
import { groupNotes } from 'web/libs/notes';
import {
  getFiltersFromSearchStr,
  Filters,
  checkFilterEqual
} from 'jslib/helpers/filters';
import { getHomePath } from 'web/libs/paths';
import NoteGroupList from './NoteGroup/List';
import HeadData from './HeadData';
import { useDispatch, useSelector } from '../../store';
import { getNotes } from '../../store/notes';
import PageToolbar from '../Common/PageToolbar';
import Paginator from '../Common/PageToolbar/Paginator';
import Flash from '../Common/Flash';
import PayWall from '../Common/PayWall';
import styles from './Home.scss';

// PER_PAGE is the number of results per page in the response from the backend implementation's API.
// Currently it is fixed.
const PER_PAGE = 30;

interface Props extends RouteComponentProps {}

function useFetchNotes(filters: Filters) {
  const dispatch = useDispatch();
  const { user } = useSelector(state => {
    return {
      user: state.auth.user.data,
      notes: state.notes
    };
  });
  const prevFilters = usePrevious(filters);

  useEffect(() => {
    if (prevFilters && checkFilterEqual(filters, prevFilters)) {
      return () => null;
    }

    dispatch(getNotes(filters));

    return () => null;
  }, [dispatch, filters, prevFilters, user]);
}

const Home: React.FunctionComponent<Props> = ({ location }) => {
  const { notes } = useSelector(state => {
    return {
      notes: state.notes
    };
  });

  const filters = getFiltersFromSearchStr(location.search);
  useFetchNotes(filters);

  const groups = groupNotes(notes.data);

  return (
    <div
      id="T-home-page"
      className="container mobile-nopadding page page-mobile-full"
    >
      <HeadData filters={filters} />

      <h1 className="sr-only">Notes</h1>

      <Flash kind="danger" when={Boolean(notes.errorMessage)}>
        Error getting notes: {notes.errorMessage}
      </Flash>

      <PageToolbar wrapperClassName={styles.toolbar}>
        <Paginator
          perPage={PER_PAGE}
          total={notes.total}
          currentPage={filters.page}
          getPath={(page: number) => {
            return getHomePath({
              ...filters.queries,
              page
            });
          }}
        />
      </PageToolbar>

      <PayWall>
        <NoteGroupList
          groups={groups}
          filters={filters}
          isFetched={notes.isFetched}
        />
      </PayWall>
    </div>
  );
};

export default withRouter(Home);
