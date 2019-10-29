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

import { RepetitionRuleData } from 'jslib/operations/types';
import RepetitionItem from './RepetitionItem';
import styles from './RepetitionList.scss';

interface Props {
  isFetching: boolean;
  isFetched: boolean;
  items: RepetitionRuleData[];
  setRuleUUIDToDelete: React.Dispatch<any>;
}

const ReptitionList: React.FunctionComponent<Props> = ({
  isFetching,
  isFetched,
  items,
  setRuleUUIDToDelete
}) => {
  if (isFetching && !isFetched) {
    return <div>loading</div>;
  }

  return (
    <ul
      id="T-repetition-rules-list"
      className={classnames('list-unstyled', styles.wrapper, {
        loaded: isFetched
      })}
    >
      {items.map(i => {
        return (
          <RepetitionItem
            key={i.uuid}
            item={i}
            setRuleUUIDToDelete={setRuleUUIDToDelete}
          />
        );
      })}
    </ul>
  );
};

export default ReptitionList;
