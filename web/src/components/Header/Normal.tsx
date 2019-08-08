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
import classnames from 'classnames';

import LogoWithText from '../Icons/LogoWithText';
import Logo from '../Icons/Logo';
import { getHomePath } from '../../libs/paths';
import AccountMenu from './AccountMenu';
import Nav from './Nav';
import SearchBar from './SearchBar';
import styles from './Normal.scss';

interface Props {}

const NormalHeader: React.SFC<Props> = () => {
  return (
    <header className={styles.wrapper}>
      <div className={classnames(styles.content, 'container mobile-nopadding')}>
        <div className={classnames(styles.left)}>
          <Link to={getHomePath({})} className={styles.brand}>
            <LogoWithText width={75} fill="white" className={styles.logo} />
            <Logo width={24} fill="white" className={styles.logosm} />
          </Link>

          <Nav />
        </div>

        <div className={classnames(styles.right)}>
          <SearchBar />

          <AccountMenu />
        </div>
      </div>
    </header>
  );
};

export default NormalHeader;
