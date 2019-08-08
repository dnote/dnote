import React from 'react';

import { monthNumToFullName } from '../../../helpers/time';
import styles from './Header.scss';

interface Props {
  year: number;
  month: number;
  total: number;
  isReady: boolean;
}

function toDatetime(year: number, month: number) {
  let monthStr;
  if (month < 10) {
    monthStr = `0${month}`;
  } else {
    monthStr = String(month);
  }

  return `${year}-${monthStr}`;
}

const Header: React.SFC<Props> = ({ year, month, total, isReady }) => {
  const monthName = monthNumToFullName(month);
  const datetime = toDatetime(year, month);

  return (
    <div className={styles.wrapper}>
      <h2 className={styles.date}>
        <time dateTime={datetime}>
          {monthName} {year}
        </time>
      </h2>
      <div className={styles.count}>{isReady && total} total</div>
    </div>
  );
};

export default Header;
