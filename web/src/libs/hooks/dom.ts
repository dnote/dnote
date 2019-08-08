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

import { useEffect } from 'react';

import { useEventListener } from './index';
import {
  KEYCODE_DOWN,
  KEYCODE_UP,
  KEYCODE_ENTER,
  KEYCODE_ESC
} from '../../helpers/keyboard';

export function useScrollToSelected({
  shouldScroll,
  offset,
  selectedOptEl,
  containerEl
}) {
  useEffect(() => {
    if (!shouldScroll || !selectedOptEl || !containerEl) {
      return;
    }

    const { scrollTop } = containerEl;
    const scrollBottom = scrollTop + containerEl.offsetHeight;
    const optionTop = selectedOptEl.offsetTop;
    const optionBottom = optionTop + selectedOptEl.offsetHeight;

    if (scrollTop > optionTop || scrollBottom < optionBottom) {
      // eslint-disable-next-line no-param-reassign
      containerEl.scrollTop = selectedOptEl.offsetTop - offset;
    }
  }, [shouldScroll, selectedOptEl, offset, containerEl]);
}

export function useScrollToFocused({
  shouldScroll,
  offset,
  focusedOptEl,
  containerEl
}) {
  useEffect(() => {
    if (!shouldScroll || !focusedOptEl || !containerEl) {
      return;
    }

    const scrollTop =
      focusedOptEl.offsetTop +
      focusedOptEl.clientHeight -
      containerEl.offsetHeight / 2 -
      offset;

    // eslint-disable-next-line no-param-reassign
    containerEl.scrollTop = scrollTop;
  }, [containerEl, focusedOptEl, offset, shouldScroll]);
}

export function useSearchMenuKeydown({
  options,
  containerEl,
  focusedIdx,
  setFocusedIdx,
  setIsOpen,
  onKeydownSelect,
  disabled
}) {
  useEventListener(containerEl, 'keydown', e => {
    if (disabled) {
      return;
    }

    const { keyCode } = e;

    if (keyCode === KEYCODE_UP || keyCode === KEYCODE_DOWN) {
      e.preventDefault();

      let nextOptionIdx;
      if (focusedIdx === 0 && keyCode === KEYCODE_UP) {
        nextOptionIdx = options.length - 1;
      } else if (
        focusedIdx === options.length - 1 &&
        keyCode === KEYCODE_DOWN
      ) {
        nextOptionIdx = 0;
      } else if (keyCode === KEYCODE_DOWN) {
        nextOptionIdx = focusedIdx + 1;
      } else if (keyCode === KEYCODE_UP) {
        nextOptionIdx = focusedIdx - 1;
      }

      setFocusedIdx(nextOptionIdx);
    } else if (keyCode === KEYCODE_ENTER) {
      e.preventDefault();
      const focusedOption = options[focusedIdx];

      if (setIsOpen) {
        setIsOpen(false);
      }

      if (onKeydownSelect) {
        onKeydownSelect(focusedOption);
      }
    } else if (keyCode === KEYCODE_ESC) {
      e.preventDefault();

      if (setIsOpen) {
        setIsOpen(false);
      }
    }
  });
}
