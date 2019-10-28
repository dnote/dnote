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
import classnames from 'classnames';

import { KEYCODE_UP, KEYCODE_DOWN } from 'jslib/helpers/keyboard';
import { useEventListener } from 'web/libs/hooks';
import Popover from '../Popover';
import { Direction, Alignment } from './types';
import styles from './Menu.scss';

export interface MenuOption {
  name: string;
  value: React.ReactElement;
}

interface ContentProps {
  options: MenuOption[];
  menuId: string;
  setContentEl: (any) => void;
  headerContent?: React.ReactNode;
}

const Content: React.FunctionComponent<ContentProps> = ({
  options,
  menuId,
  setContentEl,
  headerContent
}) => {
  return (
    <Fragment>
      {headerContent}
      <ul
        id={menuId}
        className="list-unstyled"
        role="menu"
        ref={el => {
          setContentEl(el);
        }}
      >
        {options.map(option => {
          return (
            <li role="none" key={option.name}>
              {option.value}
            </li>
          );
        })}
      </ul>
    </Fragment>
  );
};

interface MenuProps {
  options: MenuOption[];
  isOpen: boolean;
  setIsOpen: (boolean) => void;
  optRefs: React.MutableRefObject<any>[];
  triggerContent: React.ReactNode;
  triggerClassName?: string;
  contentClassName: string;
  alignment: Alignment;
  direction: Direction;
  headerContent?: React.ReactNode;
  wrapperClassName?: string;
  menuId: string;
  triggerId: string;
  disabled?: boolean;
}

const Menu: React.FunctionComponent<MenuProps> = ({
  options,
  isOpen,
  setIsOpen,
  optRefs,
  triggerContent,
  triggerClassName,
  contentClassName,
  alignment,
  direction,
  headerContent,
  wrapperClassName,
  menuId,
  triggerId,
  disabled
}) => {
  const [currentOptionIdx, setCurrentOptionIdx] = useState(0);
  const [contentEl, setContentEl] = useState(null);

  useEffect(() => {
    if (isOpen) {
      const ref = optRefs[currentOptionIdx];
      const el = ref.current;

      if (el) {
        el.focus();
      }
    } else {
      setCurrentOptionIdx(0);
    }
  }, [isOpen, currentOptionIdx, optRefs]);

  useEventListener(contentEl, 'keydown', e => {
    const { keyCode } = e;

    if (keyCode === KEYCODE_UP || keyCode === KEYCODE_DOWN) {
      // Avoid scrolling the whole page down
      e.preventDefault();
      // Stop event propagation in case any parent is also listening on the same set of keys.
      e.stopPropagation();

      let nextOptionIdx;
      if (currentOptionIdx === 0 && keyCode === KEYCODE_UP) {
        nextOptionIdx = options.length - 1;
      } else if (
        currentOptionIdx === options.length - 1 &&
        keyCode === KEYCODE_DOWN
      ) {
        nextOptionIdx = 0;
      } else if (keyCode === KEYCODE_DOWN) {
        nextOptionIdx = currentOptionIdx + 1;
      } else if (keyCode === KEYCODE_UP) {
        nextOptionIdx = currentOptionIdx - 1;
      }

      setCurrentOptionIdx(nextOptionIdx);
    }
  });

  let ariaExpanded;
  if (isOpen) {
    ariaExpanded = 'true';
  }

  return (
    <Popover
      renderTrigger={triggerProps => {
        return (
          <button
            id={triggerId}
            ref={triggerProps.triggerRef}
            type="button"
            className={classnames(
              'button button-no-ui',
              triggerProps.triggerClassName,
              triggerClassName
            )}
            onClick={() => {
              setIsOpen(!isOpen);
            }}
            aria-haspopup="menu"
            aria-expanded={ariaExpanded}
            aria-controls={menuId}
            disabled={disabled}
          >
            {triggerContent}
          </button>
        );
      }}
      contentClassName={classnames(styles.content, contentClassName)}
      wrapperClassName={wrapperClassName}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      alignment={alignment}
      direction={direction}
      renderContent={() => {
        return (
          <Content
            menuId={menuId}
            options={options}
            setContentEl={setContentEl}
            headerContent={headerContent}
          />
        );
      }}
    />
  );
};

export default Menu;
