import React from 'react';
import classnames from 'classnames';

import Item from './Item';
import { getNewPath, getBooksPath } from 'web/libs/paths';
import { Filters, toSearchObj } from 'jslib/helpers/filters';
import styles from './Nav.scss';

interface Props {
  filters: Filters;
}

const Nav: React.SFC<Props> = ({ filters }) => {
  const searchObj = toSearchObj(filters);

  return (
    <nav className={styles.wrapper}>
      <ul className={classnames('list-unstyled', styles.list)}>
        <Item to={getNewPath(searchObj)} label="New" />
        <Item to={getBooksPath(searchObj)} label="Books" />
        {/* <Item to={getRandomPath(searchObj)} label="Random" /> */}
      </ul>
    </nav>
  );
};

export default Nav;
