import React from 'react';
import classnames from 'classnames';

import styles from './DigestHolder.module.scss';

function DigestHolder() {
  return (
    <li className={styles.wrapper}>
      <div className={classnames('holder', styles.title)} />
    </li>
  );
}

export default DigestHolder;
