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

import React, { useState } from 'react';
import ReactDOM from 'react-dom';
import classnames from 'classnames';

import { Alignment, Direction } from '../Popover/types';
import styles from './Tooltip.scss';

interface Props {
  id: string;
  isOpen: boolean;
  children: React.ReactElement;
  triggerEl: HTMLElement;
  alignment: Alignment;
  direction: Direction;
}

// cumulativeOffset calculates the top and left offsets of the given element
// while taking into account all of the parents' offsets, if any.
function cumulativeOffset(element: HTMLElement) {
  let top = 0;
  let left = 0;

  let e = element;
  while (e) {
    top += e.offsetTop || 0;
    left += e.offsetLeft || 0;

    e = e.offsetParent as HTMLElement;
  }

  return {
    top,
    left
  };
}

function calcOverlayY(
  offsetY: number,
  triggerRect: ClientRect,
  overlayRect: ClientRect,
  arrowRect: ClientRect,
  alignment: Alignment,
  direction: Direction
): number {
  const triggerHeight = triggerRect.height;
  const overlayHeight = overlayRect.height;
  const arrowHeight = arrowRect.height / 2;

  if (direction === 'bottom') {
    return offsetY + triggerHeight + arrowHeight;
  }
  if (direction === 'top') {
    return offsetY - overlayHeight - arrowHeight;
  }
  if (alignment === 'bottom') {
    return offsetY + (triggerHeight - overlayHeight);
  }
  if (alignment === 'center') {
    return offsetY + (triggerHeight - overlayHeight) / 2;
  }
  if (alignment === 'top') {
    return offsetY;
  }

  return 0;
}

function calcOverlayX(
  offsetX: number,
  triggerRect: ClientRect,
  overlayRect: ClientRect,
  arrowRect: ClientRect,
  alignment: Alignment,
  direction: Direction
): number {
  const triggerWidth = triggerRect.width;
  const overlayWidth = overlayRect.width;
  const arrowWidth = arrowRect.width / 2;

  if (direction === 'left') {
    return offsetX - overlayWidth - arrowWidth;
  }
  if (direction === 'right') {
    return offsetX + triggerWidth + arrowWidth * 2;
  }
  if (alignment === 'left') {
    return offsetX;
  }
  if (alignment === 'right') {
    return offsetX + triggerWidth - overlayWidth;
  }
  if (alignment === 'center') {
    return offsetX + (triggerWidth - overlayWidth) / 2;
  }

  return 0;
}

function calcOverlayPosition(
  triggerEl: HTMLElement,
  overlayEl: HTMLElement,
  arrowEl: HTMLElement,
  direction: Direction,
  alignment: Alignment
): { top: number; left: number } {
  if (triggerEl === null) {
    return null;
  }
  if (overlayEl === null) {
    return { top: -999, left: -999 };
  }

  const triggerOffset = cumulativeOffset(triggerEl);
  const triggerRect = triggerEl.getBoundingClientRect();
  const overlayRect = overlayEl.getBoundingClientRect();
  const arrowRect = arrowEl.getBoundingClientRect();

  const x = calcOverlayX(
    triggerOffset.left,
    triggerRect,
    overlayRect,
    arrowRect,
    alignment,
    direction
  );
  const y = calcOverlayY(
    triggerOffset.top,
    triggerRect,
    overlayRect,
    arrowRect,
    alignment,
    direction
  );

  return { top: y, left: x };
}

function calcArrowX(
  offsetX: number,
  triggerRect: ClientRect,
  arrowRect: ClientRect,
  direction: Direction
) {
  const arrowWidth = arrowRect.width / 2;

  if (direction === 'top' || direction === 'bottom') {
    return offsetX + triggerRect.width / 2;
  }
  if (direction === 'left') {
    return offsetX - arrowWidth;
  }
  if (direction === 'right') {
    return offsetX + triggerRect.width;
  }

  return 0;
}

function calcArrowY(
  offsetY: number,
  triggerRect: ClientRect,
  arrowRect: ClientRect,
  direction: Direction
) {
  const arrowHeight = arrowRect.height / 2;

  if (direction === 'left' || direction === 'right') {
    return offsetY + triggerRect.height / 2 - arrowRect.height / 2;
  }
  if (direction === 'top') {
    return offsetY - arrowRect.height / 2;
  }
  if (direction === 'bottom') {
    return offsetY + triggerRect.height - arrowHeight;
  }

  return 0;
}

function calcArrowPosition(
  triggerEl: HTMLElement,
  arrowEl: HTMLElement,
  direction: Direction
) {
  if (triggerEl === null) {
    return null;
  }
  if (arrowEl === null) {
    return { top: -999, left: -999 };
  }

  const triggerOffset = cumulativeOffset(triggerEl);
  const triggerRect = triggerEl.getBoundingClientRect();
  const arrowRect = arrowEl.getBoundingClientRect();

  const x = calcArrowX(triggerOffset.left, triggerRect, arrowRect, direction);
  const y = calcArrowY(triggerOffset.top, triggerRect, arrowRect, direction);

  return { top: y, left: x };
}

const Overlay: React.FunctionComponent<Props> = ({
  id,
  isOpen,
  children,
  triggerEl,
  alignment,
  direction
}) => {
  const [overlayEl, setOverlayEl] = useState(null);
  const [arrowEl, setArrowEl] = useState(null);

  if (!isOpen) {
    return null;
  }

  const overlayRoot = document.getElementById('overlay-root');
  const overlayPos = calcOverlayPosition(
    triggerEl,
    overlayEl,
    arrowEl,
    direction,
    alignment
  );
  const arrowPos = calcArrowPosition(triggerEl, arrowEl, direction);

  return ReactDOM.createPortal(
    <div role="tooltip" id={id}>
      <div
        className={classnames(styles.arrow, {
          [styles.top]: direction === 'top',
          [styles.bottom]: direction === 'bottom',
          [styles.left]: direction === 'left',
          [styles.right]: direction === 'right'
        })}
        style={{ top: arrowPos.top, left: arrowPos.left }}
        ref={el => {
          setArrowEl(el);
        }}
      />
      <div
        className={styles.overlay}
        style={{ top: overlayPos.top, left: overlayPos.left }}
        ref={el => {
          setOverlayEl(el);
        }}
      >
        {children}
      </div>
    </div>,
    overlayRoot
  );
};

export default Overlay;
