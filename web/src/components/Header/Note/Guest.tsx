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
import { Link } from 'react-router-dom';

import Logo from '../../Icons/Logo';
import { getHomePath } from 'web/libs/paths';
import styles from './Guest.scss';

const UserNoteHeader: React.SFC = () => {
  return (
    <header className={styles.wrapper}>
      <div className={styles.content}>
        <Link to={getHomePath({})} className={styles.brand}>
          <Logo width={32} height={32} fill="#909090" className="logo" />
          <span className={styles['brand-name']}>Dnote</span>
        </Link>

        <Link
          to={getHomePath({})}
          className="button button-normal button-slim button-first-outline"
        >
          Go to Dnote
        </Link>
      </div>
    </header>
  );
};

export default UserNoteHeader;
