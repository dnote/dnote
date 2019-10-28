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

import { mdBreakpoint } from '../components/App/_variables.scss';

// catchBlur focuses rootEl if the next focused element is outside the rootEl
export function catchBlur(event, rootEl) {
  if (!rootEl) {
    return;
  }

  // If the next focus was outside the content
  if (!rootEl.contains(event.relatedTarget)) {
    rootEl.focus();
  }
}

// getScrollbarWidth measures the width of the browser's scroll bar in pixels and returns it
export function getScrollbarWidth() {
  const scrollDiv = document.createElement('div');
  scrollDiv.className = 'scrollbar-measure';
  document.body.appendChild(scrollDiv);
  const scrollbarWidth =
    scrollDiv.getBoundingClientRect().width - scrollDiv.clientWidth;
  document.body.removeChild(scrollDiv);
  return scrollbarWidth;
}

// scrollTo scrolls the given element to the given position
export function scrollTo(element: HTMLElement, posY: number) {
  if (document.body === element) {
    window.scrollTo(0, posY);
  } else {
    // eslint-disable-next-line no-param-reassign
    element.scrollTop = posY;
  }
}

// focusTextarea focuses on the given text area element and moves the cursor to the last position
export function focusTextarea(el: HTMLTextAreaElement) {
  el.focus();

  // Move the cursor to the last position
  const len = el.value.length;
  el.setSelectionRange(len, len);
}

export function checkVerticalScoll() {
  if (window.innerHeight) {
    return document.body.offsetHeight > window.innerHeight;
  }

  return (
    document.documentElement.scrollHeight >
      document.documentElement.offsetHeight ||
    document.body.scrollHeight > document.body.offsetHeight
  );
}

// getViewportDimensions returns the dimension of the viewport
export function getViewportDimensions() {
  const width = Math.max(
    document.documentElement.clientWidth,
    window.innerWidth || 0
  );
  const height = Math.max(
    document.documentElement.clientHeight,
    window.innerHeight || 0
  );

  return {
    width,
    height
  };
}

function pxToNumber(px: string): number {
  const str = px.substring(0, px.length - 2);

  return Number.parseFloat(str);
}

export function isMobileWidth() {
  const { width } = getViewportDimensions();

  const mdThreshold = pxToNumber(mdBreakpoint);

  return width < mdThreshold;
}
