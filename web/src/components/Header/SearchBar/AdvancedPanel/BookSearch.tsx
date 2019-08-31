import React, { useState, useRef, useEffect } from 'react';
import classnames from 'classnames';

import { useSelector } from '../../../../store';
import PopoverContent from '../../../Common/Popover/PopoverContent';
import { booksToOptions, filterOptions, Option } from '../../../../libs/select';
import { usePrevious } from '../../../../libs/hooks';
import styles from './AdvancedPanel.scss';

import {
  useScrollToFocused,
  useSearchMenuKeydown
} from '../../../../libs/hooks/dom';

interface Props {
  value: string;
  setValue: (string) => void;
  disabled: boolean;
}

// getCurrentTerm returns the current term in the comma separated
// value from the text input
function getCurrentTerm(searchValue: string): string {
  const parts = searchValue.split(',');
  const last = parts[parts.length - 1];

  return last.trim();
}

// getNewValue returns a new comma separated input string by appending the
// given label to the current value
function getNewValue(currentValue: string, label: string): string {
  let ret = '';

  const parts = currentValue.split(',');
  for (let i = 0; i < parts.length - 1; i++) {
    const p = parts[i].trim();

    if (p !== '') {
      ret += `${p}, `;
    }
  }

  ret += `${label},`;

  return ret;
}

// suggestionActiveRegex is the regex that matches the input value for which
// book name suggestion should be active
const suggestionActiveRegex = /.*,((?!,).)+$/;

function shouldSuggestOptions(val: string): boolean {
  if (val.trim() === '') {
    return true;
  }

  return suggestionActiveRegex.test(val);
}

function useFilteredOptions(inputValue: string) {
  const { books } = useSelector(state => {
    return {
      books: state.books.data
    };
  });

  const options = booksToOptions(books);
  const term = getCurrentTerm(inputValue);

  return filterOptions(options, term, false);
}

function useSetSuggestionVisibility(
  inputValue: string,
  setIsOpen: (boolean) => void,
  triggerRef: React.MutableRefObject<any>
) {
  const prevInputValue = usePrevious(inputValue);

  useEffect(() => {
    const triggerEl = triggerRef.current;

    if (
      shouldSuggestOptions(inputValue) &&
      document.activeElement === triggerEl
    ) {
      setIsOpen(true);
    } else if (/.*,$/.test(inputValue)) {
      setIsOpen(false);

      // focus on the input and move the cursor to the end
      if (triggerEl) {
        triggerEl.focus();
        triggerEl.selectionStart = triggerEl.value.length;
        triggerEl.selectionEnd = triggerEl.value.length;
      }
    }
  }, [setIsOpen, triggerRef, inputValue, prevInputValue]);
}

const BookSearch: React.SFC<Props> = ({ value, setValue, disabled }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [focusedIdx, setFocusedIdx] = useState(0);
  const [focusedOptEl, setFocusedOptEl] = useState(null);

  const wrapperRef = useRef(null);
  const triggerRef = useRef(null);
  const listRef = useRef(null);

  const filteredOptions = useFilteredOptions(value);

  function appendBook(o: Option) {
    const newVal = getNewValue(value, o.label);

    setFocusedIdx(0);
    setValue(newVal);
  }

  useSearchMenuKeydown({
    options: filteredOptions,
    containerEl: wrapperRef.current,
    focusedIdx,
    setFocusedIdx,
    onKeydownSelect: appendBook,
    disabled: !isOpen || disabled
  });
  useScrollToFocused({
    shouldScroll: true,
    focusedOptEl,
    containerEl: listRef.current
  });
  useSetSuggestionVisibility(value, setIsOpen, triggerRef);

  return (
    <section
      className={classnames(styles.section, styles['book-search--wrapper'])}
      ref={wrapperRef}
    >
      <label htmlFor="in-book" className={styles.label}>
        Books
        <input
          autoComplete="off"
          type="text"
          id="in-book"
          ref={triggerRef}
          className={classnames(
            'text-input text-input-small text-input-stretch',
            styles.input
          )}
          value={value}
          disabled={disabled}
          onChange={e => {
            const val = e.target.value;

            setValue(val);
          }}
          onFocus={() => {
            if (shouldSuggestOptions(value)) {
              setIsOpen(true);
            }
          }}
        />
      </label>

      <PopoverContent
        contentId="advanced-search-panel-book-list"
        onDismiss={() => {
          setIsOpen(false);
        }}
        alignment="left"
        direction="bottom"
        triggerEl={triggerRef.current}
        wrapperEl={wrapperRef.current}
        contentClassName={classnames(styles['book-suggestion-wrapper'], {
          [styles['book-suggestions-wrapper-shown']]: isOpen
        })}
        closeOnOutsideClick
        closeOnEscapeKeydown
      >
        <ul
          className={classnames(styles['book-suggestion'], 'list-unstyled')}
          ref={listRef}
        >
          {filteredOptions.map((o, idx) => {
            const isFocused = idx === focusedIdx;

            return (
              <li
                key={o.value}
                className={classnames(styles['book-item'], {})}
                ref={el => {
                  if (isFocused) {
                    setFocusedOptEl(el);
                  }
                }}
              >
                <button
                  type="button"
                  className={classnames(
                    'button-no-ui',
                    styles['book-item-button'],
                    {
                      [styles['book-item-focused']]: isFocused
                    }
                  )}
                  onClick={() => {
                    appendBook(o);
                  }}
                >
                  {o.label}
                </button>
              </li>
            );
          })}
        </ul>
      </PopoverContent>
    </section>
  );
};

export default BookSearch;
