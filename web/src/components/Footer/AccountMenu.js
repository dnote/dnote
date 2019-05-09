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

import React, { Fragment, useState, useRef } from 'react';
import classnames from 'classnames';
import { Link } from 'react-router-dom';

import Lock from '../Icons/Lock';
import Menu from '../Common/Menu';
import { signout } from '../../services/users';
import { settingsPath } from '../../libs/paths';

import styles from './AccountMenu.module.scss';

function handleSignout(e) {
  e.preventDefault();

  signout().then(() => {
    localStorage.removeItem('cipherKey');
    window.location.href = '/';
  });
}

const AccountMenu = ({ triggerClassName, demo, user }) => {
  const [isOpen, setIsOpen] = useState(false);
  const optRefs = [useRef(null), useRef(null), useRef(null)];

  const options = [
    {
      name: 'home',
      value: (
        <a
          role="menuitem"
          className={styles.link}
          href="/"
          onClick={() => {
            setIsOpen(false);
          }}
          innerRef={optRefs[0]}
          tabIndex="-1"
        >
          Home
        </a>
      )
    },
    {
      name: 'settings',
      value: (
        <Link
          role="menuitem"
          className={classnames(styles.link, {
            [styles.disabled]: demo
          })}
          to={settingsPath('account', { demo })}
          onClick={e => {
            if (demo) {
              e.preventDefault();
              return;
            }

            setIsOpen(false);
          }}
          innerRef={optRefs[1]}
          tabIndex="-1"
        >
          Settings
        </Link>
      )
    },
    {
      name: 'logout',
      value: (
        <form onSubmit={handleSignout}>
          <input
            role="menuitem"
            id="T-logout-button"
            type="submit"
            value="Logout"
            className={classnames('button-no-ui', styles.link, {
              [styles.disabled]: demo
            })}
            disabled={demo}
            ref={optRefs[2]}
            tabIndex="-1"
          />
        </form>
      )
    }
  ];

  return (
    <Menu
      options={options}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      optRefs={optRefs}
      triggerId="T-account-menu"
      menuId="account-menu"
      triggerContent={
        <Fragment>
          <span className={styles['account-circle']} />
          <span className={styles['account-label']}>Account</span>
        </Fragment>
      }
      headerContent={
        <Fragment>
          <div className={styles.header}>
            <div className={styles['session-notice-wrapper']}>
              <Lock width={16} height={16} />
              <div className={styles['session-notice']}>Signed in</div>
            </div>
            <div className={styles.email}>
              {demo ? 'user@example.com' : user.email}
            </div>
          </div>
          <div className={styles.divider} />
        </Fragment>
      }
      wrapperClassName={styles.wrapper}
      triggerClassName={triggerClassName}
      contentClassName={styles.content}
      alignment="right"
      direction="top"
    />
  );
};

export default AccountMenu;
