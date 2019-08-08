import React from 'react';
import { Link } from 'react-router-dom';
import { Location } from 'history';

import styles from './Item.scss';

interface Props {
  to: Location<any>;
  label: string;
}

const Item: React.SFC<Props> = ({ to, label }) => {
  return (
    <li className={styles.wrapper}>
      <Link to={to} className={styles.link}>
        {label}
      </Link>
    </li>
  );
};

export default Item;
