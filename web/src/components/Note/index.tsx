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

import HeaderData from './HeaderData';
import NoteContent from './NoteContent';
import Flash from '../Common/Flash';
import { getNote } from '../../store/note';
import Placeholder from './Placeholder';
import { useDispatch, useSelector, ReduxDispatch } from '../../store';
import { unsetMessage } from '../../store/ui';
import { notePathDef } from '../../libs/paths';
import { parseSearchString } from '../../libs/url';
import styles from './index.scss';

interface Match {
  noteUUID: string;
}

interface Props extends RouteComponentProps<Match> {}

function useClearMessage(dispatch: ReduxDispatch) {
  useEffect(() => {
    return () => {
      dispatch(unsetMessage(notePathDef));
    };
  }, [dispatch]);
}

function useFetchData(
  dispatch: ReduxDispatch,
  noteUUID: string,
  search: string
) {
  useEffect(() => {
    const searchObj = parseSearchString(search);

    dispatch(
      getNote(noteUUID, {
        q: searchObj.q || ''
      })
    );
  }, [dispatch, noteUUID, search]);
}

const Note: React.SFC<Props> = ({ match, location }) => {
  const { params } = match;
  const { noteUUID } = params;

  const dispatch = useDispatch();
  const { note } = useSelector(state => {
    return {
      note: state.note
    };
  });

  useFetchData(dispatch, noteUUID, location.search);
  useClearMessage(dispatch);

  if (note.errorMessage) {
    return <Flash kind="danger">Error: {note.errorMessage}</Flash>;
  }

  return (
    <div className={styles.wrapper}>
      <HeaderData note={note} />

      <div className="container mobile-nopadding page">
        {note.isFetched ? <NoteContent /> : <Placeholder />}
      </div>
    </div>
  );
};

export default React.memo(withRouter(Note));
