import React from 'react';
import classnames from 'classnames';

import { DigestData } from 'jslib/operations/types';
import Item from './Item';
import styles from './List.scss';

interface Props {
  isFetched: boolean;
  items: DigestData[];
}

const List: React.SFC<Props> = ({ items, isFetched }) => {
  if (!isFetched) {
    return <div>Loading digests...</div>;
  }

  return (
    <ul className={classnames('list-unstyled', styles.wrapper)}>
      {items.map((item, idx) => {
        const isFirst = idx === 0;
        const isLast = idx === items.length - 1;

        return (
          <Item key={item.uuid} item={item} isFirst={isFirst} isLast={isLast} />
        );
      })}
    </ul>
  );
};

export default List;
