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

import React, { useRef } from 'react';
import classnames from 'classnames';

import PopoverContent from './PopoverContent';
import styles from './Popover.module.scss';
import { Alignment, Direction } from './types';

interface Props {
  contentClassName: string;
  wrapperClassName: string;
  renderContent: () => any;
  isOpen: boolean;
  setIsOpen: (boolean) => void;
  alignment: Alignment;
  direction: Direction;
  renderTrigger: (any) => React.ReactNode;
  contentId?: string;
  closeOnOutsideClick?: boolean;
  closeOnEscapeKeydown?: boolean;
  contentHasBorder?: boolean;
  hasArrow?: boolean;
}

const Popover: React.SFC<Props> = ({
  contentClassName,
  wrapperClassName,
  renderContent,
  isOpen,
  setIsOpen,
  alignment,
  direction,
  renderTrigger,
  contentId,
  closeOnOutsideClick,
  closeOnEscapeKeydown,
  contentHasBorder,
  hasArrow
}) => {
  const triggerRef = useRef(null);

  return (
    <div className={classnames(styles.wrapper, wrapperClassName)}>
      {renderTrigger({
        triggerClassName: classnames({
          [styles['is-open']]: isOpen,
          [styles['has-arrow']]: hasArrow,
          [styles.bottom]: direction === 'bottom',
          [styles.top]: direction === 'top',
          [styles.left]: direction === 'left',
          [styles.right]: direction === 'right'
        }),
        triggerRef
      })}

      {isOpen && (
        <PopoverContent
          onDismiss={() => {
            setIsOpen(false);
          }}
          contentClassName={contentClassName}
          triggerEl={triggerRef.current}
          alignment={alignment}
          direction={direction}
          contentId={contentId}
          hasBorder={contentHasBorder}
          closeOnOutsideClick={closeOnOutsideClick}
          closeOnEscapeKeydown={closeOnEscapeKeydown}
        >
          {renderContent()}
        </PopoverContent>
      )}
    </div>
  );
};

Popover.defaultProps = {
  closeOnOutsideClick: true,
  closeOnEscapeKeydown: true,
  contentHasBorder: true,
  hasArrow: false
};

export default Popover;
