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

const tabbableNode = /input|select|textarea|button|object/;

function hidesContents(element) {
  const zeroSize = element.offsetWidth <= 0 && element.offsetHeight <= 0;

  // If the node is empty, this is good enough
  if (zeroSize && !element.innerHTML) return true;

  // Otherwise we need to check some styles
  const style = window.getComputedStyle(element);
  return zeroSize
    ? style.getPropertyValue('overflow') !== 'visible'
    : style.getPropertyValue('display') === 'none';
}

function visible(element) {
  let parentElement = element;
  while (parentElement) {
    if (parentElement === document.body) break;
    if (hidesContents(parentElement)) return false;
    parentElement = parentElement.parentNode;
  }
  return true;
}

function focusable(element, isTabIndexNotNaN) {
  const nodeName = element.nodeName.toLowerCase();

  if (!visible(element)) {
    return false;
  }

  if (nodeName === 'a') {
    return Boolean(element.href) || isTabIndexNotNaN;
  }

  return (tabbableNode.test(nodeName) && !element.disabled) || isTabIndexNotNaN;
}

function findTabbable(element) {
  return [].slice.call(element.querySelectorAll('*'), 0).filter(elm => {
    let tabIndex = elm.getAttribute('tabindex');
    if (tabIndex === null) {
      tabIndex = undefined;
    }

    const isTabIndexNaN = Number.isNaN(Number(tabIndex));

    if (!focusable(elm, !isTabIndexNaN)) {
      return false;
    }

    return isTabIndexNaN || tabIndex >= 0;
  });
}

export default function scopeTab(node, event) {
  const tabbable = findTabbable(node);

  if (!tabbable.length) {
    // Do nothing, since there are no elements that can receive focus.
    event.preventDefault();
    return;
  }

  const { shiftKey } = event;
  const head = tabbable[0];
  const tail = tabbable[tabbable.length - 1];
  let target;

  // proceed with default browser behavior on tab.
  // Focus on last element on shift + tab.
  if (node === document.activeElement) {
    if (!shiftKey) {
      return;
    }
    target = tail;
  }

  if (tail === document.activeElement && !shiftKey) {
    target = head;
  }

  if (head === document.activeElement && shiftKey) {
    target = tail;
  }

  if (target) {
    event.preventDefault();
    target.focus();
    return;
  }

  // Safari radio issue.
  //
  // Safari does not move the focus to the radio button,
  // so we need to force it to really walk through all elements.
  //
  // This is very error prune, since we are trying to guess
  // if it is a safari browser from the first occurence between
  // chrome or safari.
  //
  // The chrome user agent contains the first ocurrence
  // as the 'chrome/version' and later the 'safari/version'.
  const checkSafari = /(\bChrome\b|\bSafari\b)\//.exec(navigator.userAgent);
  const isSafariDesktop =
    checkSafari != null &&
    checkSafari[1] !== 'Chrome' &&
    /\biPod\b|\biPad\b/g.exec(navigator.userAgent) == null;

  // If we are not in safari desktop, let the browser control
  // the focus
  if (!isSafariDesktop) {
    return;
  }

  let x = tabbable.indexOf(document.activeElement);

  if (x > -1) {
    x += shiftKey ? -1 : 1;
  }

  event.preventDefault();

  tabbable[x].focus();
}
