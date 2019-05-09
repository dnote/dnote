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

import React, { useState } from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';

import SearchIcon from '../../Icons/Search';

import styles from './SearchInput.module.scss';

function SearchInput({
  value,
  setValue,
  wrapperClassName,
  inputClassName,
  placeholder,
  onFocus,
  onBlur,
  autoFocus,
  disabled,
  inputRef,
  size
}) {
  const [isFocused, setIsFocused] = useState(false);

  let iconDimension;
  if (size === 'medium') {
    iconDimension = 20;
  } else {
    iconDimension = 16;
  }

  return (
    <div
      className={classnames(styles.wrapper, wrapperClassName, {
        [styles.focused]: isFocused
      })}
    >
      <SearchIcon
        width={iconDimension}
        height={iconDimension}
        className={styles['search-icon']}
      />
      <input
        ref={inputRef}
        type="text"
        value={value}
        placeholder={placeholder}
        disabled={disabled}
        className={classnames(
          'text-input text-input-transparent',
          styles.input,
          inputClassName,
          {
            'text-input-medium': size === 'medium'
          }
        )}
        onChange={e => {
          const val = e.target.value;

          setValue(val);
        }}
        onFocus={e => {
          setIsFocused(true);

          if (onFocus) {
            onFocus(e);
          }
        }}
        onBlur={e => {
          setIsFocused(false);

          if (onBlur) {
            onBlur(e);
          }
        }}
        autoFocus={autoFocus}
      />
    </div>
  );
}

const validSizes = ['regular', 'medium'];

SearchInput.propTypes = {
  size: PropTypes.oneOf(validSizes).isRequired
};

export default SearchInput;
