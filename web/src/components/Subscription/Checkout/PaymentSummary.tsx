import React from 'react';

import styles from './Sidebar.scss';

interface Props {
  yearly: boolean;
}

const PaymentSummary: React.SFC<Props> = ({ yearly }) => {
  let total;
  if (yearly) {
    total = '$48.00';
  } else {
    total = '$5.00';
  }

  let breakdown;
  if (yearly) {
    breakdown = '$4.00 x 12 months';
  } else {
    breakdown = '';
  }

  return (
    <div className={styles['summary-wrapper']}>
      <h2 className={styles['summary-heading']}>Payment summary</h2>

      <div className={styles['summary-row']}>
        <div className={styles['summary-breakdown']}>{breakdown}</div>
        <div className={styles['summary-amount']}>{total}</div>
      </div>
    </div>
  );
};

export default PaymentSummary;
