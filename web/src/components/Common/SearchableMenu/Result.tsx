import React, { Fragment } from 'react';

import { Option } from '../../../libs/select';
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
