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

import styles from './Placeholder.scss';
import noteStyles from './NoteContent.scss';

interface Props {}

const Placeholder: React.SFC<Props> = () => {
  return (
    <div className={noteStyles.frame}>
      <div className={noteStyles.header}>
        <div className={classnames('holder', styles.title)} />
      </div>
      <div className={noteStyles.content}>
        <div className={classnames('holder', styles.line1, styles.line)} />
        <div className={classnames('holder', styles.line2, styles.line)} />
        <div className={classnames('holder', styles.line3, styles.line)} />
        <div className={classnames('holder', styles.line4, styles.line)} />
        <div className={classnames('holder', styles.line5, styles.line)} />
      </div>
    </div>
  );
};

export default Placeholder;
