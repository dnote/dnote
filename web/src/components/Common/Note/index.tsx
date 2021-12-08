/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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

import { NoteData } from 'jslib/operations/types';
import Content from './Content';
import Footer from './Footer';
import styles from './Note.scss';

interface Props {
  note: NoteData;
  header: React.ReactElement;
  footerActions?: React.ReactElement;
  headerRight?: React.ReactElement;
  collapsed?: boolean;
  footerUseTimeAgo?: boolean;
}

const Note: React.FunctionComponent<Props> = ({
  note,
  footerActions,
  collapsed,
  header,
  footerUseTimeAgo = false
}) => {
  return (
    <article
      className={classnames(styles.frame, {
        [styles.collapsed]: collapsed
      })}
    >
      {header}

      <Content note={note} collapsed={collapsed} />

      <Footer
        note={note}
        collapsed={collapsed}
        actions={footerActions}
        useTimeAgo={footerUseTimeAgo}
      />
    </article>
  );
};

export default React.memo(Note);
