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

import { Filters } from 'jslib/helpers/filters';
import { NotesGroupData } from 'web/libs/notes';
import NoteGroup from './index';
import Placeholder from './Placeholder';
import styles from './List.scss';

function renderResult({
  groups,
  isFetched,
  filters
}: {
  groups: NotesGroupData[];
  isFetched: boolean;
  filters: Filters;
}) {
  if (!isFetched) {
    return <Placeholder />;
  }

  if (groups.length === 0) {
    return <div className={styles.content}>No notes found.</div>;
  }

  return groups.map((group, idx) => {
    const isFirst = idx === 0;

    return (
      <NoteGroup
        key={`${group.year}${group.month}`}
        group={group}
        isFirst={isFirst}
        filters={filters}
      />
    );
  });
}

interface Props {
  isFetched: boolean;
  groups: NotesGroupData[];
  pro: boolean;
  filters: Filters;
}

const NoteGroupList: React.SFC<Props> = ({
  groups,
  pro,
  filters,
  isFetched
}) => {
  return (
    <div className={classnames(styles.wrapper, { [styles.nopro]: !pro })}>
      {renderResult({ groups, isFetched, filters })}
    </div>
  );
};

export default NoteGroupList;
