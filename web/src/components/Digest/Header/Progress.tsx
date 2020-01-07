/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
import classnames from 'classnames';

import { pluralize } from 'web/libs/string';
import styles from './Progress.scss';

interface Props {
  total: number;
  current: number;
}

function calcPercentage(current: number, total: number): number {
  if (total === 0) {
    return 100;
  }

  return (current / total) * 100;
}

function getCaption(current, total): string {
  if (current === total && total !== 0) {
    return 'Review completed';
  }

  return `${current} of ${total} ${pluralize('note', current)} reviewed`;
}

const Progress: React.FunctionComponent<Props> = ({ total, current }) => {
  const isComplete = current === total;
  const perc = calcPercentage(current, total);
  const width = `${perc}%`;

  return (
    <div className={styles.wrapper}>
      <div
        className={classnames(styles.caption, {
          [styles['caption-strong']]: isComplete
        })}
      >
        {getCaption(current, total)}{' '}
        <span className={styles.perc}>({perc.toFixed(0)}%)</span>
      </div>
      <div
        className={styles['bar-wrapper']}
        role="progressbar"
        aria-valuenow={perc}
        aria-valuemin={0}
        aria-valuemax={100}
      >
        <div className={styles.bar} style={{ width }} />
      </div>
    </div>
  );
};

export default Progress;
