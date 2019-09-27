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

import { getHomePath } from 'web/libs/paths';
import { useSelector } from '../../store';
import Logo from '../Icons/LogoWithText';
import styles from './SubscriptionHeader.scss';

interface Props {}

const SubscriptionsHeader: React.SFC<Props> = () => {
  const { user } = useSelector(state => {
    return {
      user: state.auth.user.data
    };
  });

  return (
    <header className={styles.wrapper}>
      <div className={styles.content}>
        <Link to={getHomePath({})} className={styles.brand}>
          <Logo
            id="subscription-header-logo"
            width={88}
            fill="black"
            className={styles.logo}
          />
        </Link>

        <div className={styles.email}>{user.email}</div>
      </div>
    </header>
  );
};

export default SubscriptionsHeader;
