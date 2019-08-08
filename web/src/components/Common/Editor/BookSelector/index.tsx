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

import React, { useState, useEffect } from 'react';
import classnames from 'classnames';

import Popover from '../../Popover';
import SearchableMenu from '../../SearchableMenu';
import BookIcon from '../../../Icons/Book';
import CaretIcon from '../../../Icons/Caret';
import SearchInput from '../../SearchInput';
import { useDispatch, useSelector } from '../../../../store';
import { BookData } from '../../../../operations/books';
import { updateBook } from '../../../../store/editor';
import { Option } from '../../../../libs/select';
import OptionItem from './OptionItem';
import styles from './index.scss';

function makeOptions(books: BookData[]): Option[] {
  const ret = [];

  for (let i = 0; i < books.length; ++i) {
    const book = books[i];

    ret.push({
      label: book.label,
      value: book.uuid
    });
  }

  return ret;
}

interface Props {
  wrapperClassName?: string;
  triggerClassName?: string;
  isReady: boolean;
  defaultIsOpen?: boolean;
  onAfterChange: () => void;
  isOpen: boolean;
  setIsOpen: (boolean) => void;
}

const BookSelector: React.SFC<Props> = ({
  wrapperClassName,
  triggerClassName,
  isReady,
  onAfterChange,
  isOpen,
  setIsOpen
}) => {
  const { books, editor } = useSelector(state => {
    return {
      books: state.books,
      editor: state.editor
    };
  });
  const dispatch = useDispatch();

  const [textboxValue, setTextboxValue] = useState('');

  useEffect(() => {
    if (!isOpen) {
      setTextboxValue('');
    }
  }, [isOpen]);

  const options = makeOptions(books.data);
  const currentValue = editor.bookUUID;
  const currentLabel = editor.bookLabel;

  let ariaExpanded;
  if (isOpen) {
    ariaExpanded = 'true';
  }

  function handleSelect(option) {
    dispatch(updateBook({ label: option.label, uuid: option.value }));
    onAfterChange();
  }

  return (
    <Popover
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      alignment="left"
      direction="bottom"
      contentClassName={classnames(styles.content)}
      wrapperClassName={classnames(styles['popover-wrapper'], wrapperClassName)}
      renderTrigger={triggerProps => {
        return (
          <button
            id="T-move-book-btn"
            ref={triggerProps.triggerRef}
            type="button"
            className={classnames(
              styles.trigger,
              triggerClassName,
              triggerProps.triggerClassName,
              {
                [styles['trigger-hidden']]: !isReady
              }
            )}
            onClick={() => {
              setIsOpen(!isOpen);
            }}
            aria-haspopup="menu"
            aria-expanded={ariaExpanded}
            aria-controls="book-filter"
            disabled={books.isFetching}
          >
            <span className={styles['book-selector-trigger']}>
              <span className={styles['book-selector-trigger-left']}>
                <BookIcon width={12} height={12} />
                <span
                  className={classnames(styles['book-label'], {
                    [styles['book-label-visible']]: Boolean(currentLabel)
                  })}
                >
                  {currentLabel || 'Choose a book'}
                </span>
              </span>
              <CaretIcon
                width={8}
                height={8}
                fill="black"
                className={styles.caret}
              />
            </span>
          </button>
        );
      }}
      renderContent={() => {
        return (
          <SearchableMenu
            menuId="book-selector"
            isOpen={isOpen}
            setIsOpen={setIsOpen}
            options={options}
            label="Choose a book"
            currentValue={currentValue}
            textboxValue={textboxValue}
            listboxClassName={styles.dropdown}
            itemClassName={styles['combobox-item']}
            labelClassName={styles['combobox-label']}
            textboxWrapperClassName={styles.combobox}
            textboxClassName={styles.textbox}
            onKeydownSelect={handleSelect}
            renderOption={(option, { isSelected, isFocused }) => {
              return (
                <OptionItem
                  option={option}
                  isFocused={isFocused}
                  isSelected={isSelected}
                  onSelect={handleSelect}
                />
              );
            }}
            renderCreateOption={(option, { isFocused }) => {
              return (
                <OptionItem
                  isNew
                  option={option}
                  isFocused={isFocused}
                  isSelected={false}
                  onSelect={handleSelect}
                />
              );
            }}
            renderInput={inputProps => {
              return (
                <SearchInput
                  size="regular"
                  placeholder="Find or create by name"
                  wrapperClassName={styles['input-wrapper']}
                  value={textboxValue}
                  setValue={setTextboxValue}
                  {...inputProps}
                />
              );
            }}
          />
        );
      }}
    />
  );
};

export default BookSelector;
