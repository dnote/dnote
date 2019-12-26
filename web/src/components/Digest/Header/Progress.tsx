import React from 'react';

import { pluralize } from 'web/libs/string';
import styles from './Progress.scss';

interface Props {
  total: number;
  current: number;
}

function calcPercentage(current: number, total: number): number {
  if (total === 0) {
    return 0;
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
  const perc = calcPercentage(current, total);
  const width = `${perc}%`;

  return (
    <div className={styles.wrapper}>
      <div className={styles.caption}>
        {getCaption(current, total)}{' '}
        <span className={styles.perc}>({perc.toFixed(0)}%)</span>
      </div>
      <div className={styles['bar-wrapper']}>
        <div className={styles.bar} style={{ width }} />
      </div>
    </div>
  );
};

export default Progress;
