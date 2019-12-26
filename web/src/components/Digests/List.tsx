import React from 'react';
import classnames from 'classnames';

import { DigestData } from 'jslib/operations/types';
import Item from './Item';
import Placeholder from './Placeholder';
import styles from './List.scss';

interface Props {
  isFetched: boolean;
  isFetching: boolean;
  items: DigestData[];
}

const List: React.FunctionComponent<Props> = ({
  items,
  isFetched,
  isFetching
}) => {
  if (isFetching) {
    return (
      <div className={styles.wrapper}>
        {[...Array(10)].map(() => {
          return <Placeholder />;
        })}
      </div>
    );
  }
  if (!isFetched) {
    return null;
  }

  return (
    <ul className={classnames('list-unstyled', styles.wrapper)}>
      {items.map(item => {
        return <Item key={item.uuid} item={item} />;
      })}
    </ul>
  );
};

export default List;
