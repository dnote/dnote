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

import React from 'react';
import Helmet from 'react-helmet';

import { DigestData } from 'jslib/operations/types';
import { DigestNoteData } from 'jslib/operations/types';
import { getDigestTitle } from './helpers';
import { useDispatch } from '../../store';
import { setDigestNoteReviewed } from '../../store/digest';
import Placeholder from '../Common/Note/Placeholder';
import NoteItem from './NoteItem';
import styles from './Digest.scss';

interface Props {
  notes: DigestNoteData[];
  digest: DigestData;
  isFetched: boolean;
  isFetching: boolean;
}

const NoteList: React.FunctionComponent<Props> = ({
  isFetched,
  isFetching,
  notes,
  digest
}) => {
  const dispatch = useDispatch();

  function handleSetReviewed(noteUUID: string, isReviewed: boolean) {
    return dispatch(
      setDigestNoteReviewed({ digestUUID: digest.uuid, noteUUID, isReviewed })
    );
  }

  if (isFetching) {
    return (
      <div className={styles.wrapper}>
        <Placeholder wrapperClassName={styles.item} />
        <Placeholder wrapperClassName={styles.item} />
        <Placeholder wrapperClassName={styles.item} />
      </div>
    );
  }
  if (!isFetched) {
    return null;
  }

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>{`${getDigestTitle(digest)} - Digest`}</title>
      </Helmet>
      <ul className={styles.list}>
        {notes.map(note => {
          return (
            <NoteItem
              key={note.uuid}
              note={note}
              onSetReviewed={handleSetReviewed}
            />
          );
        })}
      </ul>
    </div>
  );
};

export default NoteList;
