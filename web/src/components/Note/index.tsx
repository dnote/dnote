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

import React, { useEffect, useState } from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import { parseSearchString } from 'jslib/helpers/url';
import { checkOwner } from 'web/libs/notes';
import Content from '../Common/Note';
import Placeholder from '../Common/Note/Placeholder';
import Flash from '../Common/Flash';
import { getNote } from '../../store/note';
import { useDispatch, useSelector, ReduxDispatch } from '../../store';
import DeleteModal from './DeleteModal';
import ShareModal from './ShareModal';
import HeaderData from './HeaderData';
import FooterActions from './FooterActions';
import Header from './Header';
import styles from './index.scss';

interface Match {
  noteUUID: string;
}

interface Props extends RouteComponentProps<Match> {}

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

const Note: React.FunctionComponent<Props> = ({ match, location }) => {
  const { params } = match;
  const { noteUUID } = params;

  const dispatch = useDispatch();
  const { note, user } = useSelector(state => {
    return {
      note: state.note,
      user: state.auth.user
    };
  });
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isShareModalOpen, setIsShareModalOpen] = useState(false);

  const isOwner = checkOwner(note.data, user.data);

  useFetchData(dispatch, noteUUID, location.search);

  if (note.errorMessage) {
    return <Flash kind="danger">Error: {note.errorMessage}</Flash>;
  }

  return (
    <div id="T-note-page" className={styles.wrapper}>
      <HeaderData note={note} />

      <div className="container mobile-nopadding page page-mobile-full">
        {note.isFetched ? (
          <Content
            note={note.data}
            footerActions={
              <FooterActions
                isOwner={isOwner}
                noteUUID={note.data.uuid}
                onDeleteModalOpen={() => {
                  setIsDeleteModalOpen(true);
                }}
                onShareModalOpen={() => {
                  setIsShareModalOpen(true);
                }}
              />
            }
            header={<Header isOwner={isOwner} note={note.data} />}
          />
        ) : (
          <Placeholder />
        )}
      </div>

      <DeleteModal
        isOpen={isDeleteModalOpen}
        onDismiss={() => {
          setIsDeleteModalOpen(false);
        }}
        noteUUID={note.data.uuid}
      />

      <ShareModal
        isOpen={isShareModalOpen}
        onDismiss={() => {
          setIsShareModalOpen(false);
        }}
        note={note.data}
      />
    </div>
  );
};

export default React.memo(withRouter(Note));
