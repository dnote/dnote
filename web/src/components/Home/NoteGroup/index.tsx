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

import React, { Fragment } from 'react';
import classnames from 'classnames';

import NoteItem from '../NoteItem';
import Header from './Header';
import Flash from '../../Common/Flash';
import NoteHolder from './NoteHolder';
import { NotesGroup } from '../../../store/notes';
import styles from './NoteGroup.scss';

function renderPlaceholders() {
  return (
    <Fragment>
      <NoteHolder />
      <NoteHolder />
      <NoteHolder />
      <NoteHolder />
      <NoteHolder />
      <NoteHolder />
      <NoteHolder />
      <NoteHolder />
    </Fragment>
  );
}

function renderItems(group: NotesGroup) {
  const { uuids, items } = group;

  return uuids.map(uuid => {
    const item = items[uuid];

    return <NoteItem key={uuid} note={item} />;
  });
}

interface Props {
  group: NotesGroup;
  isFirst: boolean;
}

const NoteGroup: React.SFC<Props> = ({ group, isFirst }) => {
  const {
    year,
    month,
    error,
    total,
    isFetching,
    isFetched,
    isFetchingMore
  } = group;

  return (
    <div className={classnames(styles.wrapper, { [styles.first]: isFirst })}>
      <Header year={year} month={month} total={total} isReady={isFetched} />

      <Flash kind="danger" when={Boolean(error)}>
        Could not get notes: {error}
      </Flash>

      <ul className={classnames('list-unstyled', styles.list)}>
        {isFetching ? renderPlaceholders() : renderItems(group)}
        {isFetchingMore && renderPlaceholders()}
      </ul>
    </div>
  );
};

export default NoteGroup;
