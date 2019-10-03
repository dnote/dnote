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

import { useEffect, useRef } from 'react';

import {
  KEYCODE_DOWN,
  KEYCODE_UP,
  KEYCODE_ENTER,
  KEYCODE_ESC
} from 'jslib/helpers/keyboard';
import { Option } from 'jslib/helpers/select';
import { useEventListener } from './index';
import { scrollTo } from '../dom';

interface ScrollToSelectedParams {
  shouldScroll: boolean;
  offset: number;
  selectedOptEl: HTMLElement;
  containerEl: HTMLElement;
}

export function useScrollToSelected({
  shouldScroll,
  offset,
  selectedOptEl,
  containerEl
}: ScrollToSelectedParams) {
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

interface ScrollToFocusedParms {
  shouldScroll: boolean;
  offset?: number;
  focusedOptEl: HTMLElement;
  containerEl: HTMLElement;
}

export function useScrollToFocused({
  shouldScroll,
  offset = 0,
  focusedOptEl,
  containerEl
}: ScrollToFocusedParms) {
  useEffect(() => {
    if (!shouldScroll || !focusedOptEl || !containerEl) {
      return;
    }

    let visibleHeight;
    if (containerEl === document.body) {
      visibleHeight = window.innerHeight;
    } else {
      visibleHeight = containerEl.offsetHeight;
    }

    const posY =
      focusedOptEl.offsetTop +
      focusedOptEl.clientHeight -
      visibleHeight / 2 -
      offset;

    scrollTo(containerEl, posY);
  }, [containerEl, focusedOptEl, offset, shouldScroll]);
}

export type KeydownSelectFn<T> = (T) => void;

interface SearchMenuKeydownParams<T> {
  options: T[];
  containerEl: HTMLElement;
  focusedIdx: number;
  setFocusedIdx: (number) => void;
  onKeydownSelect: KeydownSelectFn<T>;
  setIsOpen?: (boolean) => void;
  disabled?: boolean;
}

export function useSearchMenuKeydown<T = Option>({
  options,
  containerEl,
  focusedIdx,
  setFocusedIdx,
  setIsOpen,
  onKeydownSelect,
  disabled
}: SearchMenuKeydownParams<T>) {
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
      e.stopPropagation();

      if (setIsOpen) {
        setIsOpen(false);
      }
    }
  });
}

export function useFocus() {
  const elRef = useRef<HTMLElement>();

  const setFocus = () => {
    const currentEl = elRef.current;

    if (currentEl) {
      currentEl.focus();
    }
  };

  return [setFocus, elRef] as const;
}
