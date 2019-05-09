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
import { connect } from 'react-redux';
import classnames from 'classnames';

import AccountMenu from './AccountMenu';

import styles from './Footer.module.scss';

const Footer = ({ user, demo }) => {
  return (
    <footer className={styles.footer}>
      <ul className={classnames('list-unstyled', styles['action-list'])}>
        <li>
          <AccountMenu
            tirggerClassName={styles.action}
            user={user}
            demo={demo}
          />
        </li>
      </ul>
    </footer>
  );
};

function mapStateToProps(state) {
  return {
    user: state.auth.user.data
  };
}

export default connect(mapStateToProps)(Footer);
