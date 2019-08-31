import React from 'react';
import { Link } from 'react-router-dom';

import { SettingSections, getSettingsPath } from '../../libs/paths';
import styles from './MobileMenu.scss';

interface Props {}

const MobileMenu: React.SFC<Props> = () => {
  return (
    <div className={styles.wrapper}>
      <ul className="list-unstyled">
        <li>
          <Link
            className={styles.link}
            to={getSettingsPath(SettingSections.account)}
          >
            Settings
          </Link>
        </li>
        <li>
          <Link className={styles.link} to="/s">
            Logout
          </Link>
        </li>
      </ul>
    </div>
  );
};

export default MobileMenu;
