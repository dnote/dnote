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

import React, { useState, useRef, useEffect } from 'react';
import { withRouter } from 'react-router-dom';
import classnames from 'classnames';

import styles from './SearchableMenu.module.scss';
import { makeOptionId, getOptIdxByValue } from '../../../helpers/accessibility';
import {
  useScrollToFocused,
  useScrollToSelected,
  useSearchMenuKeydown
} from '../../../libs/hooks/dom';

function useFocusedIdx(options, currentValue) {
  const initialValue = getOptIdxByValue(options, currentValue);

  return useState(initialValue);
}

function filterOptions(options, textboxValue) {
  if (!textboxValue) {
    return options;
  }

  const searchReg = new RegExp(`${textboxValue}`, 'i');
  return options.filter(book => {
    // The built-in 'All notes' option has a value of an empty string
    return searchReg.test(book.label) && book.value !== '';
  });
}

function getScrollOffset(headerEl) {
  let ret = 0;
  if (headerEl) {
    ret = headerEl.offsetHeight;
  }

  return ret;
}

function SearchableMenu({
  menuId,
  isOpen,
  setIsOpen,
  options,
  label,
  listboxClassName,
  currentValue,
  textboxValue,
  itemClassName,
  labelClassName,
  textboxWrapperClassName,
  textboxClassName,
  renderOption,
  history,
  onKeydownSelect,
  location,
  match,
  renderInput,
  disabled
}) {
  const [containerEl, setContainerEl] = useState(null);
  const [wrapperEl, setWrapperEl] = useState(null);
  const [focusedOptEl, setFocusedOptEl] = useState(null);
  const headerRef = useRef(null);
  const selectedOptRef = useRef(null);

  const filteredOptions = filterOptions(options, textboxValue);
  const [focusedIdx, setFocusedIdx] = useFocusedIdx(
    filteredOptions,
    currentValue
  );

  const currentIdx = getOptIdxByValue(filteredOptions, currentValue);
  useEffect(() => {
    setFocusedIdx(currentIdx);
  }, [textboxValue, setFocusedIdx, currentIdx]);

  const offset = getScrollOffset(headerRef.current);
  useScrollToSelected({
    shouldScroll: isOpen,
    offset,
    containerEl,
    selectedOptEl: selectedOptRef.current
  });
  useScrollToFocused({
    shouldScroll: isOpen,
    focusedOptEl,
    containerEl,
    offset
  });
  useSearchMenuKeydown({
    options: filteredOptions,
    containerEl: wrapperEl,
    focusedIdx,
    setFocusedIdx,
    setIsOpen,
    onKeydownSelect,
    location,
    match,
    history,
    disabled
  });

  const currentOptId = makeOptionId(menuId, currentValue);

  return (
    <div
      ref={el => {
        setWrapperEl(el);
      }}
    >
      <header ref={headerRef}>
        <span className={labelClassName}>{label}</span>

        <div className={textboxWrapperClassName}>
          {renderInput({
            autoFocus: true,
            type: 'text',
            value: textboxValue,
            placeholder: 'Find your book',
            className: classnames(styles['search-input'], textboxClassName)
          })}
        </div>
      </header>

      <ul
        ref={el => {
          setContainerEl(el);
        }}
        id={menuId}
        className={classnames('list-unstyled', listboxClassName)}
        tabIndex="0"
        role="menu"
        aria-activedescendant={currentOptId}
      >
        {filteredOptions.map((option, idx) => {
          const id = makeOptionId(menuId, option.value);

          const isSelected = option.value === currentValue;
          const isFocused = idx === focusedIdx;

          return (
            // eslint-disable-next-line jsx-a11y/click-events-have-key-events
            <li
              id={id}
              key={option.value}
              className={classnames(itemClassName)}
              role="none"
              onClick={() => {
                if (disabled) {
                  return;
                }

                setIsOpen(false);
              }}
              ref={el => {
                if (isSelected) {
                  selectedOptRef.current = el;
                }

                if (isFocused) {
                  setFocusedOptEl(el);
                }
              }}
            >
              {renderOption(option, { isSelected, isFocused })}
            </li>
          );
        })}
      </ul>
    </div>
  );
}

export default withRouter(SearchableMenu);
