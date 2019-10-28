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
import classnames from 'classnames';

import { booksToOptions, filterOptions, Option } from 'jslib/helpers/select';
import { KEYCODE_BACKSPACE } from 'jslib/helpers/keyboard';
import { useSearchMenuKeydown, useScrollToFocused } from 'web/libs/hooks/dom';
import { useSelector } from '../../store';
import PopoverContent from '../Common/Popover/PopoverContent';
import CloseIcon from '../Icons/Close';
import { usePrevious } from 'web/libs/hooks';
import styles from './MultiSelect.scss';

function getTextInputWidth(term: string, active: boolean) {
  if (!active && term === '') {
    return '100%';
  }

  const val = 14 + term.length * 8;
  return `${val}px`;
}

interface Props {
  options: Option[];
  currentOptions: Option[];
  setCurrentOptions: (Option) => void;
  placeholder?: string;
  disabled?: boolean;
  textInputId?: string;
  wrapperClassName?: string;
  inputInnerRef?: React.MutableRefObject<any>;
}

// TODO: Make a generic Select component that works for both single and multiple selection
// by passing of a flag
const MultiSelect: React.SFC<Props> = ({
  options,
  currentOptions,
  setCurrentOptions,
  textInputId,
  placeholder,
  disabled,
  wrapperClassName,
  inputInnerRef
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [focusedIdx, setFocusedIdx] = useState(0);
  const [focusedOptEl, setFocusedOptEl] = useState(null);
  const [term, setTerm] = useState('');

  const wrapperRef = useRef(null);
  const inputRef = useRef(null);
  const listRef = useRef(null);

  const currentValues = currentOptions.map(o => {
    return o.value;
  });
  const possibleOptions = options.filter(o => {
    return currentValues.indexOf(o.value) === -1;
  });

  const filteredOptions = filterOptions(possibleOptions, term, false);

  function appendOption(o: Option | undefined) {
    if (!o) {
      return;
    }

    setTerm('');
    const newVal = [...currentOptions, o];
    setCurrentOptions(newVal);
  }
  function removeOption(o: Option) {
    setTerm('');
    const newVal = currentOptions.filter(opt => {
      return opt.value !== o.value;
    });
    setCurrentOptions(newVal);
  }
  function popOption() {
    if (currentOptions.length === 0) {
      return;
    }

    const newVal = currentOptions.slice(0, -1);
    setCurrentOptions(newVal);
  }

  useSearchMenuKeydown({
    options: filteredOptions,
    containerEl: wrapperRef.current,
    focusedIdx,
    setFocusedIdx,
    onKeydownSelect: appendOption,
    disabled: !isOpen || disabled
  });
  useScrollToFocused({
    shouldScroll: true,
    focusedOptEl,
    containerEl: listRef.current
  });

  useEffect(() => {
    if (!isOpen) {
      inputRef.current.blur();
      setTerm('');
    }
  }, [isOpen]);

  // useEffect(() => {
  //   if (term !== '' && !isOpen) {
  //     setIsOpen(true);
  //   }
  // }, [term, isOpen]);

  const active = currentOptions.length > 0;
  const textInputWidth = getTextInputWidth(term, active);

  return (
    <div
      className={classnames('form-select', styles.wrapper, wrapperClassName, {
        [styles.disabled]: disabled,
        'form-select-disabled': disabled
      })}
      ref={wrapperRef}
      onClick={() => {
        if (inputRef.current) {
          inputRef.current.focus();
        }

        // setIsOpen(!isOpen);
      }}
    >
      <ul className={styles['current-options']}>
        <span
          className={classnames(styles.placeholder, {
            [styles.hidden]: active || term !== ''
          })}
        >
          {placeholder}
        </span>
        {currentOptions.map(o => {
          return (
            <li className={styles['current-option-item']} key={o.value}>
              <div className={styles['current-option-label']}>{o.label}</div>
              <button
                type="button"
                className={classnames('button-no-ui', styles['dismiss-option'])}
                aria-label="Remove the option"
                onClick={e => {
                  if (!isOpen) {
                    e.stopPropagation();
                    if (inputRef.current) {
                      inputRef.current.focus();
                    }
                  }
                  removeOption(o);
                }}
              >
                <CloseIcon width={12} height={12} />
              </button>
            </li>
          );
        })}
        <li className={styles['input-wrapper']}>
          <input
            autoComplete="off"
            type="text"
            id={textInputId}
            ref={el => {
              inputRef.current = el;

              if (inputInnerRef) {
                inputInnerRef.current = el;
              }
            }}
            className={classnames(styles.input, {
              [styles['active-input']]: active
            })}
            value={term}
            disabled={disabled}
            onKeyDown={e => {
              if (e.keyCode === KEYCODE_BACKSPACE) {
                if (term === '') {
                  popOption();
                }
              }
            }}
            onChange={e => {
              const val = e.target.value;

              setTerm(val);
            }}
            onClick={e => {
              e.preventDefault();
            }}
            onFocus={() => {
              setIsOpen(true);
            }}
            style={{
              width: textInputWidth
            }}
          />
        </li>
      </ul>

      <PopoverContent
        contentId="advanced-search-panel-book-list"
        onDismiss={() => {
          setIsOpen(false);
        }}
        alignment="left"
        direction="bottom"
        triggerEl={inputRef.current}
        wrapperEl={wrapperRef.current}
        contentClassName={classnames(styles['suggestion-wrapper'], {
          [styles['suggestions-wrapper-shown']]: isOpen
        })}
        closeOnOutsideClick
        closeOnEscapeKeydown
      >
        <ul
          className={classnames(styles['suggestion'], 'list-unstyled')}
          ref={listRef}
        >
          {filteredOptions.map((o, idx) => {
            const isFocused = idx === focusedIdx;

            return (
              <li
                key={o.value}
                className={classnames(styles['suggestion-item'], {})}
                ref={el => {
                  if (isFocused) {
                    setFocusedOptEl(el);
                  }
                }}
              >
                <button
                  type="button"
                  tabIndex={-1}
                  className={classnames(
                    'button-no-ui',
                    styles['suggestion-item-button'],
                    {
                      [styles['suggestion-item-focused']]: isFocused
                    }
                  )}
                  onClick={() => {
                    appendOption(o);
                  }}
                >
                  {o.label}
                </button>
              </li>
            );
          })}
        </ul>
      </PopoverContent>
    </div>
  );
};

export default MultiSelect;
