import React from 'react';

import { monthNumToFullName } from '../../../helpers/time';
import styles from './Header.scss';

interface Props {
  year: number;
  month: number;
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

const Header: React.SFC<Props> = ({ year, month }) => {
  const monthName = monthNumToFullName(month);
  const datetime = toDatetime(year, month);

  return (
    <header className={styles.wrapper}>
      <h2 className={styles.date}>
        <time dateTime={datetime}>
          {monthName} {year}
        </time>
      </h2>
    </header>
  );
};

export default Header;
