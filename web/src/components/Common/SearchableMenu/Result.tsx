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

import React, { Fragment } from 'react';

import { Option } from 'jslib/helpers/select';
import { makeOptionId } from '../../../helpers/accessibility';
import Item from './Item';

interface Props {
  options: Option[];
  menuId: string;
  currentValue: string;
  focusedIdx: number;
  disabled: boolean;
  itemClassName: string;
  setIsOpen: (boolean) => void;
  selectedOptRef: React.MutableRefObject<any>;
  setFocusedOptEl: (HTMLElement) => void;
  renderOption: (Option, OptionParams) => React.ReactNode;
  renderCreateOption: (Option, OptionParams) => React.ReactNode;
}

const Result: React.SFC<Props> = ({
  options,
  menuId,
  currentValue,
  focusedIdx,
  disabled,
  itemClassName,
  setIsOpen,
  selectedOptRef,
  renderOption,
  renderCreateOption,
  setFocusedOptEl
}) => {
  return (
    <Fragment>
      {options.map((option, idx) => {
        const id = makeOptionId(menuId, option.value);

        const isSelected = option.value === currentValue;
        const isFocused = idx === focusedIdx;

        return (
          // eslint-disable-next-line jsx-a11y/click-events-have-key-events
          <Item
            id={id}
            key={option.value}
            itemClassName={itemClassName}
            disabled={disabled}
            value={option.value}
            setIsOpen={setIsOpen}
            selectedOptRef={selectedOptRef}
            setFocusedOptEl={setFocusedOptEl}
            isSelected={isSelected}
            isFocused={isFocused}
          >
            {option.value === ''
              ? renderCreateOption(option, { isFocused })
              : renderOption(option, { isSelected, isFocused })}
          </Item>
        );
      })}
    </Fragment>
  );
};

export default Result;
