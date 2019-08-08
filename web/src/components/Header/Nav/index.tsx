import React from 'react';
import classnames from 'classnames';

import Item from './Item';
import { getNewPath, getBooksPath, getRandomPath } from '../../../libs/paths';
import styles from './Nav.scss';

interface Props {}

const Nav: React.SFC<Props> = () => {
  return (
    <nav className={styles.wrapper}>
      <ul className={classnames('list-unstyled', styles.list)}>
        <Item to={getNewPath()} label="New" />
        <Item to={getBooksPath()} label="Books" />
        <Item to={getRandomPath()} label="Random" />
      </ul>
    </nav>
  );
};

export default Nav;
