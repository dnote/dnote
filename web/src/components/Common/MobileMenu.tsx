import React from 'react';
import classnames from 'classnames';
import { Link } from 'react-router-dom';

import CloseIcon from '../Icons/Close';
import { SettingSections, getSettingsPath } from '../../libs/paths';
import * as usersService from '../../services/users';
import styles from './MobileMenu.scss';

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
            <form
              onSubmit={e => {
                e.preventDefault();

                usersService.signout().then(() => {
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
