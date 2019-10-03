import React from 'react';
import classnames from 'classnames';

import { RepetitionRuleData } from 'jslib/operations/types';
import RepetitionItem from './RepetitionItem';
import styles from './RepetitionList.scss';

interface Props {
  isFetching: boolean;
  isFetched: boolean;
  items: RepetitionRuleData[];
}

const ReptitionList: React.SFC<Props> = ({ items }) => {
  return (
    <ul className={classnames('list-unstyled', styles.wrapper)}>
      {items.map(i => {
        return <RepetitionItem key={i.uuid} item={i} />;
      })}
    </ul>
  );
};

export default ReptitionList;
