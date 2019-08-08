import React from 'react';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import styles from './Item.scss';

interface Props {
  to: string;
  label: string;
  renderIcon: (string) => React.ReactNode;
  active: boolean;
}

const Item: React.SFC<Props> = ({ to, label, renderIcon, active }) => {
  let fill;
  if (active) {
    fill = '#49abfd';
  } else {
    fill = '#cecece';
  }

  return (
    <li className={styles.wrapper}>
      <Link
        to={to}
        className={classnames(styles.link, { [styles.active]: active })}
      >
        {renderIcon(fill)}
        <div className={styles.label}>{label}</div>
      </Link>
    </li>
  );
};

export default Item;
