import React from 'react';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import { DigestData } from 'jslib/operations/types';
import { getDigestPath } from 'web/libs/paths';
import Time from '../Common/Time';
import { timeAgo } from '../../helpers/time';
import styles from './Item.scss';

interface Props {
  item: DigestData;
}

const Item: React.FunctionComponent<Props> = ({ item }) => {
  const createdAt = new Date(item.createdAt);

  return (
    <li
      className={classnames(styles.wrapper, {
        [styles.read]: item.isRead,
        [styles.unread]: !item.isRead
      })}
    >
      <Link to={getDigestPath(item.uuid)} className={styles.link}>
        <span className={styles.title}>
          {item.repetitionRule.title} #{item.version}
        </span>
        <Time
          id={`${item.uuid}-ts`}
          text={timeAgo(createdAt.getTime())}
          ms={createdAt.getTime()}
          wrapperClassName={styles.ts}
        />
      </Link>
    </li>
  );
};

export default Item;
