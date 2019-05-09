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
import moment from 'moment';
import classnames from 'classnames';

import NoteItem from './NoteItem';
import NoteHolder from './NoteHolder';
import { presentError } from '../../helpers/error';
import styles from './NoteGroup.module.scss';

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

function renderItems({ items, demo, currentNoteUUID, search }) {
  return Object.keys(items).map(uuid => {
    const item = items[uuid];
    const { data, errorMessage } = item;

    const isSelected = data.uuid === currentNoteUUID;

    return (
      <NoteItem
        note={data}
        key={data.uuid}
        demo={demo}
        isSelected={isSelected}
        search={search}
        errorMessage={errorMessage}
      />
    );
  });
}

function getItems(group) {
  const { uuids } = group;

  return uuids.map(uuid => {
    return group.items[uuid];
  });
}

export default function({
  group,
  demo,
  year,
  month,
  error,
  total,
  isFetching,
  isFetched,
  isFetchingMore,
  currentNoteUUID,
  bookFilterIsOpen,
  search,
  isFirst
}) {
  const monthName = moment.months(month - 1);
  const items = getItems(group);

  return (
    <div className={classnames(styles.wrapper, { [styles.first]: isFirst })}>
      {bookFilterIsOpen && <div className={styles.mask} />}

      <div className={styles.header}>
        <div className={styles['header-date']}>
          {monthName} {year}
        </div>
        <div className={styles['header-count']}>{isFetched && total}</div>
      </div>

      {error && (
        <div className={styles.error}>
          Could not get notes: {presentError(error)}
        </div>
      )}

      <ul className={classnames('list-unstyled', styles.list)}>
        {isFetching
          ? renderPlaceholders()
          : renderItems({ items, demo, currentNoteUUID, search })}
        {isFetchingMore && renderPlaceholders()}
      </ul>
    </div>
  );
}
