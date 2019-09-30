import React from 'react';
import classnames from 'classnames';

interface Props {
  id: string;
  value: string;
  itemClassName: string;
  children: React.ReactNode;
  disabled: boolean;
  setIsOpen: (boolean) => void;
  selectedOptRef: React.MutableRefObject<any>;
  setFocusedOptEl: (HTMLElement) => void;
  isFocused: boolean;
  isSelected?: boolean;
}

const Item: React.SFC<Props> = ({
  children,
  id,
  value,
  disabled,
  itemClassName,
  setIsOpen,
  selectedOptRef,
  setFocusedOptEl,
  isSelected,
  isFocused
}) => {
  return (
    <li
      id={id}
      key={value}
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
          // eslint-disable-next-line no-param-reassign
          selectedOptRef.current = el;
        }

        if (isFocused) {
          setFocusedOptEl(el);
        }
      }}
    >
      {children}
    </li>
  );
};

export default Item;
