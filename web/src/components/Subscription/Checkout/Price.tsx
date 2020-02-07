import React from 'react';
import classnames from 'classnames';

import styles from './Sidebar.scss';

interface Props {
  yearly: boolean;
}

const Price: React.SFC<Props> = ({ yearly }) => {
  return (
    <div className={styles['price-wrapper']}>
      {yearly ? (
        <span>
          <span className={styles.price}>$4.00</span>
          <span className={classnames(styles.price, styles['price-strike'])}>
            $5.00
          </span>
        </span>
      ) : (
        <span className={styles.price}>$5.00</span>
      )}
      <div className={styles.interval}>/ month</div>
    </div>
  );
};

export default Price;
