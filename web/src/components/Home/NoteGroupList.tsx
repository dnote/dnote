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
import { withRouter, RouteComponentProps } from 'react-router-dom';

import NoteGroup from './NoteGroup';
import styles from './NoteGroupList.scss';
import { NotesGroup } from '../../store/notes';

interface Props extends RouteComponentProps {
  groups: NotesGroup[];
  pro: boolean;
}

const NoteGroupList: React.SFC<Props> = ({ groups, pro }) => {
  return (
    <div className={classnames(styles.wrapper, { [styles.nopro]: !pro })}>
      {groups.length === 0 ? (
        <div className={styles.content}>No notes found.</div>
      ) : (
        groups.map((group, idx) => {
          const isFirst = idx === 0;

          return (
            <NoteGroup
              key={`${group.year}${group.month}`}
              group={group}
              isFirst={isFirst}
            />
          );
        })
      )}
    </div>
  );
};

export default withRouter(NoteGroupList);
