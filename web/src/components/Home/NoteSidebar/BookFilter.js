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
import { Link, withRouter } from 'react-router-dom';
import classnames from 'classnames';

import SearchableMenu from '../../Common/SearchableMenu';
import SearchInput from '../../Common/SearchInput';
import Popover from '../../Common/Popover';
import CaretIcon from '../../Icons/Caret';
import CheckIcon from '../../Icons/Check';
import { getHomePath, getNotePath, isNotePath } from '../../../libs/paths';
import { parseSearchString } from '../../../libs/url';
import { getFacetsFromSearchStr } from '../../../libs/facets';

import styles from './BookFilter.module.scss';
import SearchMenuStyles from '../../Common/SearchableMenu/General.module.scss';

function getOptionByValue(options, value) {
  for (let i = 0; i < options.length; ++i) {
    const o = options[i];

    if (o.value === value) {
      return o;
    }
  }

  return {};
}

function getOptions(books) {
  const ret = [];

  ret.push({
    label: 'All notes',
    value: ''
  });

  for (let i = 0; i < books.length; ++i) {
    const book = books[i];

    ret.push({
      label: book.label,
      value: book.uuid
    });
  }

  return ret;
}

function getOptionDestination({ demo, location, match, option }) {
  let ret;

  if (isNotePath(location.pathname)) {
    const { params } = match;

    const searchObj = parseSearchString(location.search);

    const newSearchObj = {
      ...searchObj,
      book: option.value
    };

    ret = getNotePath(params.noteUUID, newSearchObj, {
      demo,
      isEditor: true
    });
  } else {
    ret = getHomePath({ book: option.value }, { demo });
  }

  return ret;
}

function handleMenuKeydownSelect({ demo, location, match, history }) {
  return option => {
    const destination = getOptionDestination({ demo, location, match, option });

    history.push(destination);
  };
}

function renderOption(
  option,
  { isSelected, isFocused },
  { demo, location, match }
) {
  const destination = getOptionDestination({
    demo,
    location,
    match,
    option
  });

  return (
    <Link
      role="menuitem"
      to={destination}
      className={classnames(SearchMenuStyles['combobox-option'], {
        [SearchMenuStyles.active]: isSelected,
        [SearchMenuStyles.focused]: isFocused
      })}
      tabIndex="-1"
    >
      {isSelected && (
        <CheckIcon
          className={SearchMenuStyles['check-icon']}
          fill="white"
          width={12}
          height={12}
        />
      )}

      <div className={styles['book-label']}>{option.label}</div>
    </Link>
  );
}

function getCurrentValue(location) {
  const facets = getFacetsFromSearchStr(location.search);

  if (!facets.book) {
    return '';
  }

  return facets.book;
}

const BookFilter = ({
  books,
  isFetching,
  noteSidebarOpen,
  isOpen,
  setIsOpen,
  demo,
  location,
  match,
  history
}) => {
  const [textboxValue, setTextboxValue] = useState('');
  const options = getOptions(books);

  const currentValue = getCurrentValue(location);
  const currentOpt = getOptionByValue(options, currentValue);

  useEffect(() => {
    if (!isOpen) {
      setTextboxValue('');
    }
  }, [isOpen]);

  useEffect(() => {
    if (!noteSidebarOpen) {
      setIsOpen(false);
    }
  }, [noteSidebarOpen, setIsOpen]);

  let ariaExpanded;
  if (isOpen) {
    ariaExpanded = 'true';
  }

  return (
    <Popover
      renderTrigger={({ triggerClassName, triggerRef }) => {
        return (
          <button
            ref={triggerRef}
            type="button"
            className={classnames(
              'button button-no-ui',
              styles.trigger,
              triggerClassName
            )}
            onClick={() => {
              setIsOpen(!isOpen);
            }}
            aria-haspopup="menu"
            aria-expanded={ariaExpanded}
            aria-controls="book-filter"
            disabled={isFetching}
          >
            <span className={styles['button-content']}>
              <div className={styles['button-label']}>{currentOpt.label}</div>
              <CaretIcon
                width={12}
                height={12}
                fill="#6e6e6e"
                className={SearchMenuStyles.caret}
              />
            </span>
          </button>
        );
      }}
      contentClassName={classnames(styles.content, SearchMenuStyles.content)}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      alignment="left"
      direction="bottom"
      renderContent={() => {
        return (
          <SearchableMenu
            menuId="book-filter"
            isOpen={isOpen}
            setIsOpen={setIsOpen}
            options={options}
            listboxClassName={SearchMenuStyles.dropdown}
            label="Filter books"
            currentValue={currentValue}
            textboxValue={textboxValue}
            itemClassName={SearchMenuStyles['combobox-item']}
            labelClassName={SearchMenuStyles['combobox-label']}
            textboxWrapperClassName={SearchMenuStyles.combobox}
            textboxClassName={SearchMenuStyles.textbox}
            onKeydownSelect={handleMenuKeydownSelect({
              demo,
              location,
              match,
              history
            })}
            renderOption={(option, { isSelected, isFocused }) => {
              return renderOption(
                option,
                { isSelected, isFocused },
                { demo, location, match }
              );
            }}
            renderInput={inputProps => {
              return (
                <SearchInput
                  size="regular"
                  placeholder="Find a book"
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

export default withRouter(BookFilter);
