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
import { Link } from 'react-router-dom';

import {
  SettingSections,
  getSettingsPath,
  getRepetitionsPath
} from 'web/libs/paths';
import services from 'web/libs/services';
import styles from './MobileMenu.scss';
import CloseIcon from '../Icons/Close';

interface Props {
  onDismiss: () => void;
  isOpen: boolean;
}

const MobileMenu: React.SFC<Props> = ({ onDismiss, isOpen }) => {
  if (!isOpen) {
    return null;
  }

  return (
    <nav className={styles.wrapper} aria-labelledby="mobile-menu">
      <div className="sr-only" id="mobile-menu">
        Mobile menu
      </div>

      <div className={styles['close-wrapper']}>
        <button
          onClick={onDismiss}
          type="button"
          aria-label="Close the modal"
          className={classnames('button-no-ui', styles.close)}
        >
          <CloseIcon width={24} height={24} fill="white" />
        </button>
      </div>

      <div className={styles.section}>
        <span className={styles.subheading}>Menu</span>

        <ul className={classnames('list-unstyled', styles.list)}>
          <li className={styles.item}>
            <Link
              className={styles.link}
              to={getSettingsPath(SettingSections.account)}
            >
              Settings
            </Link>
          </li>
          <li className={styles.item}>
            <Link className={styles.link} to={getRepetitionsPath()}>
              Repetition
            </Link>
          </li>
          <li className={styles.item}>
            <form
              onSubmit={e => {
                e.preventDefault();

                services.users.signout().then(() => {
                  window.location.href = '/';
                });
              }}
            >
              <input
                type="submit"
                value="Logout"
                className={classnames(
                  'button-no-ui',
                  styles.link,
                  styles['logout-button']
                )}
              />
            </form>
          </li>
        </ul>
      </div>
    </nav>
  );
};

export default MobileMenu;
