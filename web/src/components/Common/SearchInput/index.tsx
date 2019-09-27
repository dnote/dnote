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

import Actions from './Actions';
import styles from './SearchInput.scss';

interface Props {
  value: string;
  onChange: (Event) => void;
  placeholder: string;
  inputClassName?: string;
  wrapperClassName?: string;
  autoFocus?: boolean;
  disabled?: boolean;
  inputRef?: React.MutableRefObject<any>;
  onFocus?: (Event) => void;
  onBlur?: (Event) => void;
  onReset?: (Event) => void;
  expanded?: boolean;
  setExpanded?: (boolean) => void;
  inputId?: string;
}

const SearchInput: React.SFC<Props> = ({
  value,
  onChange,
  inputClassName,
  wrapperClassName,
  placeholder,
  onFocus,
  onBlur,
  autoFocus,
  disabled,
  inputRef,
  onReset,
  expanded,
  setExpanded,
  inputId
}: Props) => {
  const resettable = Boolean(onReset);
  const expandable = Boolean(setExpanded);

  return (
    <div className={classnames(styles.wrapper, wrapperClassName)}>
      <label htmlFor={inputId} className="sr-only">
        {placeholder}
      </label>

      <input
        id={inputId}
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
            [styles.resettable]: resettable,
            [styles.expandable]: expandable
          }
        )}
        onChange={onChange}
        onFocus={e => {
          if (onFocus) {
            onFocus(e);
          }
        }}
        onBlur={e => {
          if (onBlur) {
            onBlur(e);
          }
        }}
        autoFocus={autoFocus}
      />

      <Actions
        onReset={onReset}
        expanded={expanded}
        setExpanded={setExpanded}
        resetShown={value !== ''}
      />
    </div>
  );
};

export default SearchInput;
