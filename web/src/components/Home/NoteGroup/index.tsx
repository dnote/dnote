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

import { NotesGroupData } from 'web/libs/notes';
import { Filters } from 'jslib/helpers/filters';
import NoteItem from './NoteItem';
import Header from './Header';
import styles from './NoteGroup.scss';

function renderItems(group: NotesGroupData, filters: Filters) {
  return group.data.map(item => {
    return <NoteItem key={item.uuid} note={item} filters={filters} />;
  });
}

interface Props {
  group: NotesGroupData;
  isFirst: boolean;
  filters: Filters;
}

const NoteGroup: React.SFC<Props> = ({ group, isFirst, filters }) => {
  const { year, month } = group;

  return (
    <section
      className={classnames(styles.wrapper, { [styles.first]: isFirst })}
    >
      <Header year={year} month={month} />

      <ul className={classnames('list-unstyled', styles.list)}>
        {renderItems(group, filters)}
      </ul>
    </section>
  );
};

export default NoteGroup;
