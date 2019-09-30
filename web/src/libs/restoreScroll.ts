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

/*
Copyright (c) 2013-2015 Brigade
https://www.brigade.com/

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

/* eslint-disable */

// script restoreScroll provides a feature to restore the browser scroll
// upon navigating backwards.

if (window.history.pushState) {
  const SCROLL_RESTORATION_TIMEOUT_MS = 3000;
  const TRY_TO_SCROLL_INTERVAL_MS = 50;

  const originalPushState = window.history.pushState;
  const originalReplaceState = window.history.replaceState;

  // Store current scroll position in current state when navigating away.
  window.history.pushState = function() {
    const newStateOfCurrentPage = Object.assign({}, window.history.state, {
      __scrollX: window.scrollX,
      __scrollY: window.scrollY
    });
    originalReplaceState.call(window.history, newStateOfCurrentPage, '');

    originalPushState.apply(window.history, arguments);
  };

  // Make sure we don't throw away scroll position when calling "replaceState".
  window.history.replaceState = function(state, ...otherArgs) {
    const newState = Object.assign(
      {},
      {
        __scrollX: window.history.state && window.history.state.__scrollX,
        __scrollY: window.history.state && window.history.state.__scrollY
      },
      state
    );

    originalReplaceState.apply(window.history, [newState].concat(otherArgs));
  };

  let timeoutHandle = null;
  let scrollBarWidth = null;

  // Try to scroll to the scrollTarget, but only if we can actually scroll
  // there. Otherwise keep trying until we time out, then scroll as far as
  // we can.
  const tryToScrollTo = scrollTarget => {
    // Stop any previous calls to "tryToScrollTo".
    clearTimeout(timeoutHandle);

    const body = document.body;
    const html = document.documentElement;
    if (!scrollBarWidth) {
      scrollBarWidth = getScrollbarWidth();
    }

    // From http://stackoverflow.com/a/1147768
    const documentWidth = Math.max(
      body.scrollWidth,
      body.offsetWidth,
      html.clientWidth,
      html.scrollWidth,
      html.offsetWidth
    );
    const documentHeight = Math.max(
      body.scrollHeight,
      body.offsetHeight,
      html.clientHeight,
      html.scrollHeight,
      html.offsetHeight
    );

    if (
      (documentWidth + scrollBarWidth - window.innerWidth >= scrollTarget.x &&
        documentHeight + scrollBarWidth - window.innerHeight >=
          scrollTarget.y) ||
      Date.now() > scrollTarget.latestTimeToTry
    ) {
      window.scrollTo(scrollTarget.x, scrollTarget.y);
    } else {
      timeoutHandle = setTimeout(
        () => tryToScrollTo(scrollTarget),
        TRY_TO_SCROLL_INTERVAL_MS
      );
    }
  };

  // Try scrolling to the previous scroll position on popstate
  const onPopState = () => {
    const state = window.history.state;

    if (
      state &&
      Number.isFinite(state.__scrollX) &&
      Number.isFinite(state.__scrollY)
    ) {
      setTimeout(() =>
        tryToScrollTo({
          x: state.__scrollX,
          y: state.__scrollY,
          latestTimeToTry: Date.now() + SCROLL_RESTORATION_TIMEOUT_MS
        })
      );
    }
  };

  // Calculating width of browser's scrollbar
  function getScrollbarWidth() {
    let outer = document.createElement('div');
    outer.style.visibility = 'hidden';
    outer.style.width = '100px';
    outer.style.msOverflowStyle = 'scrollbar';

    document.body.appendChild(outer);

    let widthNoScroll = outer.offsetWidth;
    // force scrollbars
    outer.style.overflow = 'scroll';

    // add innerdiv
    let inner = document.createElement('div');
    inner.style.width = '100%';
    outer.appendChild(inner);

    let widthWithScroll = inner.offsetWidth;

    // remove divs
    outer.parentNode.removeChild(outer);

    return widthNoScroll - widthWithScroll;
  }

  window.addEventListener('popstate', onPopState, true);
}
