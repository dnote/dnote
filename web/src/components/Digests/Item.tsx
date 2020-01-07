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
      <Link
        id={`digest-item-${item.uuid}`}
        to={getDigestPath(item.uuid)}
        className={styles.link}
      >
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
