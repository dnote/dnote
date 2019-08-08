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
import { withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';

import Result from './Result';
import { makeOptionId, getOptIdxByValue } from '../../../helpers/accessibility';
import { Option } from '../../../libs/select';
import {
  useScrollToFocused,
  useScrollToSelected,
  useSearchMenuKeydown
} from '../../../libs/hooks/dom';
import styles from './SearchableMenu.scss';

function useFocusedIdx(options: Option[], currentValue) {
  const initialValue = getOptIdxByValue(options, currentValue);

  return useState(initialValue);
}

function filterOptions(
  options: Option[],
  textboxValue: string,
  creatable: boolean
): Option[] {
  if (!textboxValue) {
    return options;
  }

  const ret = [];
  const searchReg = new RegExp(`${textboxValue}`, 'i');
  let hit = null;

  for (let i = 0; i < options.length; i++) {
    const option = options[i];

    if (option.label === textboxValue) {
      hit = option;
    } else if (searchReg.test(option.label) && option.value !== '') {
      ret.push(option);
    }
  }

  // if there is an exact match, display the option at the top
  // otherwise, display a creatable option at the bottom
  if (hit) {
    ret.unshift(hit);
  } else if (creatable) {
    // creatable option has a value of an empty string
    ret.push({ label: textboxValue, value: '' });
  }

  return ret;
}

function getScrollOffset(headerEl) {
  let ret = 0;
  if (headerEl) {
    ret = headerEl.offsetHeight;
  }

  return ret;
}

interface OptionParms {
  isSelected: boolean;
  isFocused: boolean;
}

interface Props extends RouteComponentProps<any> {
  menuId: string;
  isOpen: boolean;
  setIsOpen: (boolean) => void;
  options: Option[];
  label: string;
  listboxClassName: string;
  currentValue: string;
  textboxValue: string;
  itemClassName: string;
  labelClassName: string;
  textboxWrapperClassName: string;
  textboxClassName: string;
  renderOption: (Option, OptionParams) => React.ReactNode;
  renderCreateOption: (Option, OptionParams) => React.ReactNode;
  onKeydownSelect: (Option) => void;
  renderInput: (any) => React.ReactNode;
  disabled?: boolean;
}

const SearchableMenu: React.SFC<Props> = ({
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
  onKeydownSelect,
  renderInput,
  disabled,
  renderCreateOption
}) => {
  const [containerEl, setContainerEl] = useState(null);
  const [wrapperEl, setWrapperEl] = useState(null);
  const [focusedOptEl, setFocusedOptEl] = useState(null);
  const headerRef = useRef(null);
  const selectedOptRef = useRef(null);

  const creatable = Boolean(renderCreateOption);

  const filteredOptions = filterOptions(options, textboxValue, creatable);
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
        tabIndex={0}
        role="menu"
        aria-activedescendant={currentOptId}
      >
        <Result
          options={filteredOptions}
          menuId={menuId}
          currentValue={currentValue}
          focusedIdx={focusedIdx}
          disabled={disabled}
          itemClassName={itemClassName}
          setIsOpen={setIsOpen}
          selectedOptRef={selectedOptRef}
          renderOption={renderOption}
          renderCreateOption={renderCreateOption}
          setFocusedOptEl={setFocusedOptEl}
        />
      </ul>
    </div>
  );
};

export default withRouter(SearchableMenu);
