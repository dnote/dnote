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

import React, { Fragment, useState, useEffect } from 'react';
import { connect } from 'react-redux';
import classnames from 'classnames';

import Popover from '../../Common/Popover';
import SearchableMenu from '../../Common/SearchableMenu';
import CheckIcon from '../../Icons/Check';
import BookIcon from '../../Icons/Book';
import CaretIcon from '../../Icons/Caret';
import { updateBookUUID, markDirty } from '../../../actions/editor';
import SearchInput from '../../Common/SearchInput';

import styles from './BookSelector.module.scss';
import searchMenuStyles from '../../Common/SearchableMenu/General.module.scss';

function renderOption(option, { isSelected, isFocused }, demo, handleSelect) {
  return (
    <button
      role="option"
      type="button"
      disabled={demo}
      aria-selected={isSelected.toString()}
      onClick={() => {
        handleSelect(option);
      }}
      className={classnames(
        'button-no-ui',
        `book-item-${option.value}`,
        searchMenuStyles['combobox-option'],
        {
          [searchMenuStyles.active]: isSelected,
          [searchMenuStyles.focused]: isFocused
        }
      )}
    >
      {isSelected && (
        <CheckIcon
          fill="white"
          width={12}
          height={12}
          className={searchMenuStyles['check-icon']}
        />
      )}
      <div className={styles['option-label']}>{option.label}</div>
    </button>
  );
}

function getOptions(books) {
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

function getOptionByValue(options, value) {
  for (let i = 0; i < options.length; ++i) {
    const o = options[i];

    if (o.value === value) {
      return o;
    }
  }

  return {};
}

function BookSelector({
  booksData,
  editorData,
  doUpdateBookUUID,
  doMarkDirty,
  wrapperClassName,
  triggerClassName,
  isReady,
  demo
}) {
  const [isOpen, setIsOpen] = useState(false);
  const [textboxValue, setTextboxValue] = useState('');

  useEffect(() => {
    if (!isOpen) {
      setTextboxValue('');
    }
  }, [isOpen]);

  const options = getOptions(booksData.items);
  const currentValue = editorData.bookUUID;
  const currentOption = getOptionByValue(options, currentValue);

  let ariaExpanded;
  if (isOpen) {
    ariaExpanded = 'true';
  }

  function handleSelect(option) {
    doUpdateBookUUID(option.value);
    doMarkDirty();
  }

  return (
    <Popover
      renderTrigger={triggerProps => {
        return (
          <button
            id="T-move-book-btn"
            ref={triggerProps.triggerRef}
            type="button"
            className={classnames(
              'button button-no-ui',
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
            disabled={booksData.isFetching}
          >
            <span className={styles['book-selector-trigger']}>
              <BookIcon width={12} height={12} />
              <div className={styles['book-label']}>{currentOption.label}</div>
              <CaretIcon
                width={8}
                height={8}
                fill="black"
                className={searchMenuStyles.caret}
              />
            </span>
          </button>
        );
      }}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      alignment="left"
      direction="bottom"
      contentClassName={classnames(styles.content, searchMenuStyles.content)}
      wrapperClassName={classnames(styles['popover-wrapper'], wrapperClassName)}
      renderContent={() => {
        return (
          <Fragment>
            <SearchableMenu
              identifier="book-selector"
              isOpen={isOpen}
              setIsOpen={setIsOpen}
              options={options}
              label="Move the note"
              currentValue={currentValue}
              textboxValue={textboxValue}
              setTextboxValue={setTextboxValue}
              listboxClassName={searchMenuStyles.dropdown}
              itemClassName={searchMenuStyles['combobox-item']}
              activeItemClassName={searchMenuStyles['combobox-active-item']}
              labelClassName={searchMenuStyles['combobox-label']}
              textboxWrapperClassName={searchMenuStyles.combobox}
              textboxClassName={searchMenuStyles.textbox}
              onKeydownSelect={handleSelect}
              renderOption={(option, { isSelected, isFocused }) => {
                return renderOption(
                  option,
                  { isSelected, isFocused },
                  demo,
                  handleSelect
                );
              }}
              disabled={demo}
              renderInput={inputProps => {
                return (
                  <SearchInput
                    disabled={demo}
                    size="regular"
                    placeholder="Find your book"
                    wrapperClassName={styles['input-wrapper']}
                    value={textboxValue}
                    setValue={setTextboxValue}
                    {...inputProps}
                  />
                );
              }}
            />
          </Fragment>
        );
      }}
    />
  );
}

function mapStateToProps(state) {
  return {
    booksData: state.books,
    editorData: state.editor
  };
}

const mapDispatchToProps = {
  doUpdateBookUUID: updateBookUUID,
  doMarkDirty: markDirty
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(BookSelector);
