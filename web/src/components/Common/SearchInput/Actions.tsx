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

import CloseIcon from '../../Icons/Close';
import CaretIcon from '../../Icons/CaretSolid';
import styles from './SearchInput.scss';

interface Props {
  onReset?: (Event) => void;
  resetShown?: boolean;
  expanded?: boolean;
  setExpanded?: (boolean) => void;
}

const Actions: React.SFC<Props> = ({
  onReset,
  resetShown,
  setExpanded,
  expanded
}) => {
  const resettable = Boolean(onReset);
  const expandable = Boolean(setExpanded);

  return (
    <div className={styles.actions}>
      {resettable && (
        <button
          type="button"
          className={classnames('button-no-ui', styles.reset, {
            [styles['reset-shown']]: resetShown
          })}
          aria-label="Clear search"
          onClick={onReset}
        >
          <CloseIcon width={16} height={16} />
        </button>
      )}
      {expandable && (
        <button
          type="button"
          className={classnames('button-no-ui', styles.expand)}
          aria-label="Expand search menu"
          onClick={() => {
            setExpanded(!expanded);
          }}
        >
          <CaretIcon width={12} height={12} />
        </button>
      )}
    </div>
  );
};

export default Actions;
