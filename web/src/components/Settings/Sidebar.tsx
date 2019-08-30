import React from 'react';
import classnames from 'classnames';
import { NavLink } from 'react-router-dom';

import { SettingSections, getSettingsPath } from '../../libs/paths';
import styles from './Sidebar.scss';

interface Props {}

const Sidebar: React.SFC<Props> = () => {
  return (
    <nav className={styles.wrapper}>
      <ul className={classnames('list-unstyled')}>
        <li>
          <NavLink
            className={styles.item}
            activeClassName={styles.active}
            to={getSettingsPath(SettingSections.account)}
          >
            Account
          </NavLink>
        </li>
        <li>
          <NavLink
            className={styles.item}
            activeClassName={styles.active}
            to={getSettingsPath(SettingSections.notification)}
          >
            Notification
          </NavLink>
        </li>
        <li>
          <NavLink
            className={styles.item}
            activeClassName={styles.active}
            to={getSettingsPath(SettingSections.billing)}
          >
            Billing
          </NavLink>
        </li>
      </ul>
    </nav>
  );
};

export default Sidebar;
