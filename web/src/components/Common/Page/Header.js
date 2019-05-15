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

import { toggleSidebar } from '../../../actions/ui';
import SidebarToggle from '../SidebarToggle';

import styles from './Header.module.scss';

function Header({ doToggleSidebar, heading, leftContent, rightContent }) {
  return (
    <div className={styles.header}>
      <div className="container-wide">
        <div className={styles.content}>
          <div className={styles.left}>
            <SidebarToggle onClick={doToggleSidebar} />

            <div>
              <h1 className={styles.heading}>{heading}</h1>

              {leftContent}
            </div>
          </div>

          {rightContent}
        </div>
      </div>
    </div>
  );
}

const mapDispatchToProps = {
  doToggleSidebar: toggleSidebar
};

export default connect(
  null,
  mapDispatchToProps
)(Header);
