import React from 'react';

import formatDate from 'web/helpers/time/format';
import styles from './Sidebar.scss';

interface Props {
  yearly: boolean;
}

const PaymentSummary: React.SFC<Props> = ({ yearly }) => {
  let interval;
  if (yearly) {
    interval = 'year';
  } else {
    interval = 'month';
  }

  let adj;
  if (yearly) {
    adj = 'yearly';
  } else {
    adj = 'monthly';
  }

  let cost;
  if (yearly) {
    cost = '$48.00';
  } else {
    cost = '$5.00';
  }

  const startDate = formatDate(new Date(), '%M/%D/%YYYY ');

  return (
    <p className={styles['schedule-summary']}>
      Your {adj} plan starts on {startDate}
      at {cost} per {interval} and will renew after a {interval}. Cancel
      auto-renewal any time.
    </p>
  );
};

export default PaymentSummary;
