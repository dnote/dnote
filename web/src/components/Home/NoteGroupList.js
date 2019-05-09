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
import classnames from 'classnames';

import NoteGroup from './NoteGroup';
import styles from './NoteGroupList.module.scss';

function Empty() {
  return <div className={styles.content}>No notes found.</div>;
}

function List({ groups, demo, bookFilterIsOpen, currentNoteUUID, location }) {
  return groups.map((group, idx) => {
    const isFirst = idx === 0;

    return (
      <NoteGroup
        key={`${group.year}${group.month}`}
        group={group}
        year={group.year}
        month={group.month}
        total={group.total}
        error={group.error}
        isFetching={group.isFetching}
        isFetched={group.isFetched}
        isFetchingMore={group.isFetchingMore}
        hasFetchedMore={group.hasFetchedMore}
        demo={demo}
        bookFilterIsOpen={bookFilterIsOpen}
        currentNoteUUID={currentNoteUUID}
        search={location.search}
        isFirst={isFirst}
      />
    );
  });
}

export default function({
  groups,
  demo,
  bookFilterIsOpen,
  currentNoteUUID,
  location,
  cloud
}) {
  return (
    <div className={classnames(styles.wrapper, { [styles.nopro]: !cloud })}>
      {groups.length === 0 ? (
        <Empty />
      ) : (
        <List
          groups={groups}
          demo={demo}
          bookFilterIsOpen={bookFilterIsOpen}
          currentNoteUUID={currentNoteUUID}
          location={location}
        />
      )}
    </div>
  );
}
