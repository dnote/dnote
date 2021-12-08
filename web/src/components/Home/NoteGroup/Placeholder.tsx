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

import React, { Fragment } from 'react';
import classnames from 'classnames';

import styles from './Placeholder.scss';

const NoteHolder = () => {
  return (
    <div className={styles['note-wrapper']}>
      <div className={classnames('holder', styles.book)} />
      <div className={classnames('holder', styles.ts)} />
      <div className={classnames('holder', styles.content)} />
    </div>
  );
};

const HeaderHolder = () => {
  return (
    <div className={styles['header-wrapper']}>
      <div className={classnames('holder', styles.heading)} />
    </div>
  );
};

const Placeholder: React.FunctionComponent = () => {
  return (
    <Fragment>
      <HeaderHolder />
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
};

export default Placeholder;
