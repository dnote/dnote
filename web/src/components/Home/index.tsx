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
import { Location } from 'history';
// import classnames from 'classnames';
//
// import Workspace from './Workspace';
//
// import { resetNote } from '../../actions/note';
// import { resetEditor, stageNote } from '../../actions/editor';
// import { getCipherKey } from '../../crypto';
// import { getNote } from '../../actions/note';
// import { usePrevious } from '../../libs/hooks';
//
// import style from './Home.module.scss';

import NoteGroupList from './NoteGroupList';
import HeadData from './HeadData';
import { useDispatch, useSelector } from '../../store';
import { getInitialNotes, resetNotes } from '../../store/notes';
import { getFacetsFromSearchStr } from '../../libs/facets';

interface Props extends RouteComponentProps {}

function useFetchInitialNotes(location: Location<any>) {
  const dispatch = useDispatch();
  const { user, notes } = useSelector(state => {
    return {
      user: state.auth.user.data,
      notes: state.notes
    };
  });

  useEffect(() => {
    if (notes.initialized) {
      return () => null;
    }
    if (user.uuid === '' || !user.pro) {
      return () => null;
    }

    const date = new Date();
    const year = date.getUTCFullYear();
    const month = date.getUTCMonth() + 1;
    const facets = getFacetsFromSearchStr(location.search);

    dispatch(resetNotes());
    dispatch(
      getInitialNotes({
        facets,
        year,
        month
      })
    );

    return () => null;
  }, [location.search, user, dispatch, notes.initialized]);
}

const Home: React.SFC<Props> = ({ location }) => {
  const { groups, user } = useSelector(state => {
    return {
      user: state.auth.user.data,
      groups: state.notes.groups
    };
  });

  useFetchInitialNotes(location);

  return (
    <div className="container mobile-nopadding">
      <HeadData />

      <h1 className="sr-only">Notes</h1>

      <NoteGroupList groups={groups} pro={user.pro} />
    </div>
  );
};

export default withRouter(Home);
