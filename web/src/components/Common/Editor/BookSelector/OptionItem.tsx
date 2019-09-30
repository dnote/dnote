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

import { Option } from 'jslib/helpers/select';
import CheckIcon from '../../../Icons/Check';
import styles from './OptionItem.scss';

interface Props {
  option: Option;
  isSelected: boolean;
  isFocused: boolean;
  onSelect: (Option) => void;
  isNew?: boolean;
  id?: string;
  className?: string;
}

function renderBody(isNew: boolean, label: string) {
  if (isNew) {
    return `Create book '${label}'`;
  }

  return label;
}

const OptionItem: React.SFC<Props> = ({
  option,
  isSelected,
  isFocused,
  onSelect,
  isNew,
  id
}) => {
  return (
    <button
      id={id}
      role="option"
      type="button"
      aria-selected={isSelected}
      onClick={() => {
        onSelect(option);
      }}
      className={classnames(
        'T-book-item-option',
        'button-no-ui',
        `book-item-${option.value}`,
        styles['combobox-option'],
        {
          [styles.active]: isSelected,
          [styles.focused]: isFocused
        }
      )}
    >
      {isSelected && (
        <CheckIcon
          fill="white"
          width={12}
          height={12}
          className={styles['check-icon']}
        />
      )}
      <span className={styles['option-label']}>
        {renderBody(isNew, option.label)}
      </span>
    </button>
  );
};

export default OptionItem;
