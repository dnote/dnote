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
import classnames from 'classnames';

import { DigestData } from 'jslib/operations/types';
import { getRange } from 'jslib/helpers/arr';
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
        {getRange(10).map(key => {
          return <Placeholder key={key} />;
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
