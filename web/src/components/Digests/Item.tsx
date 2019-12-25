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
  isFirst: boolean;
  isLast: boolean;
}

const Item: React.SFC<Props> = ({ item, isFirst, isLast }) => {
  const createdAt = new Date(item.createdAt);
  const isRead = item.receipts.length > 0;

  return (
    <li
      className={classnames(styles.wrapper, {
        [styles.first]: isFirst,
        [styles.last]: isLast,
        [styles.read]: isRead,
        [styles.unread]: !isRead
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
