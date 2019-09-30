/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

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
