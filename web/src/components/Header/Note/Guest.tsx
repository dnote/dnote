/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import { getHomePath } from 'web/libs/paths';
import LogoWithText from '../../Icons/LogoWithText';
import styles from './Guest.scss';

const UserNoteHeader: React.FunctionComponent = () => {
  return (
    <header className={styles.wrapper}>
      <div className={styles.content}>
        <Link to={getHomePath({})} className={styles.brand}>
          <LogoWithText id="main-logo-text" width={75} fill="#909090" />
        </Link>

        <Link to={getHomePath()} className={styles.cta}>
          Go to Dnote &#8250;
        </Link>
      </div>
    </header>
  );
};

export default UserNoteHeader;
