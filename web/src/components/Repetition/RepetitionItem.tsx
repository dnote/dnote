import React from 'react';
import classnames from 'classnames';

import { RepetitionRuleData } from 'jslib/operations/types';
import styles from './RepetitionItem.scss';

interface Props {
  item: RepetitionRuleData;
}

const RepetitionItem: React.SFC<Props> = ({ item }) => {
  return (
    <li className={styles.wrapper}>
      <div className={styles.left}>
        <div
          className={classnames(styles.status, {
            [styles.active]: item.enabled
          })}
        >
          Enabled
        </div>
      </div>

      <div>
        <h2 className={styles.title}>{item.title}</h2>
      </div>
    </li>
  );
};

export default RepetitionItem;
