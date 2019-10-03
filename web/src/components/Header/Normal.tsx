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
import { Link, withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';
import { Location } from 'history';

import { getHomePath, checkCurrentPath, homePathDef } from 'web/libs/paths';
import { toSearchObj } from 'jslib/helpers/filters';
import LogoWithText from '../Icons/LogoWithText';
import Logo from '../Icons/Logo';
import AccountMenu from './AccountMenu';
import Nav from './Nav';
import SearchBar from './SearchBar';
import { useFilters } from '../../store';
import { FiltersState } from '../../store/filters';
import styles from './Normal.scss';

interface Props extends RouteComponentProps {}

function getHomeDest(location: Location, filters: FiltersState) {
  if (checkCurrentPath(location, homePathDef)) {
    return getHomePath();
  }

  return getHomePath(filters);
}

const NormalHeader: React.SFC<Props> = ({ location }) => {
  const filters = useFilters();
  const searchObj = toSearchObj(filters);

  return (
    <header id="T-main-header" className={styles.wrapper}>
      <div
        className={classnames(styles['content-wrapper'], 'container mobile-fw')}
      >
        <div className={classnames(styles.content)}>
          <div className={classnames(styles.left)}>
            <Link
              id="T-home-link"
              to={getHomeDest(location, searchObj)}
              className={styles.brand}
            >
              <LogoWithText
                id="main-logo-text"
                width={75}
                fill="white"
                className={styles.logo}
              />
              <Logo width={24} fill="white" className={styles.logosm} />
            </Link>

            <Nav filters={filters} />
          </div>

          <div className={classnames(styles.right)}>
            <SearchBar />

            <AccountMenu />
          </div>
        </div>
      </div>
    </header>
  );
};

export default withRouter(NormalHeader);
