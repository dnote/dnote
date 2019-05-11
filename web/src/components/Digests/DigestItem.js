import React, { useState } from 'react';
import classnames from 'classnames';
import moment from 'moment';

import { Link } from 'react-router-dom';

import styles from './DigestItem.module.scss';
import { digestPath } from '../../libs/paths';

function DigestItem({ digest, demo }) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <li
      className={classnames(styles.wrapper, {
        [styles.active]: isHovered
      })}
      key={digest.uuid}
      onMouseEnter={() => {
        setIsHovered(true);
      }}
      onMouseLeave={() => {
        setIsHovered(false);
      }}
    >
      <Link className={styles.link} to={digestPath(digest.uuid, { demo })}>
        {moment(digest.created_at).format('YYYY MMM Do')}
      </Link>
    </li>
  );
}

export default DigestItem;
