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

import React, { useEffect, useCallback, useState } from 'react';
import ReactDOM from 'react-dom';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';

import { KEYCODE_ESC, KEYCODE_TAB } from 'jslib/helpers/keyboard';
import scopeTab from 'web/libs/scopeTab';
import { getScrollbarWidth } from 'web/libs/dom';
import { useEventListener } from 'web/libs/hooks';
import styles from './Modal.scss';

let scrollbarWidth = 0;

interface Props extends RouteComponentProps {
  isOpen: boolean;
  onDismiss: () => void;
  ariaLabelledBy: string;
  modalId?: string;
  ariaDescribedBy?: string;
  size?: string;
  overlayClassName?: string;
  modalClassName?: string;
}

const Modal: React.SFC<Props> = ({
  isOpen,
  onDismiss,
  overlayClassName,
  modalClassName,
  ariaLabelledBy,
  ariaDescribedBy,
  size = 'regular',
  children,
  modalId
}) => {
  const [contentEl, setContentEl] = useState(null);

  const focusContent = useCallback(() => {
    if (!contentEl) {
      return;
    }
    // check if content has focus
    if (
      document.activeElement === contentEl ||
      contentEl.contains(document.activeElement)
    ) {
      return;
    }

    contentEl.focus();
  }, [contentEl]);

  const handleMousedown = useCallback(
    e => {
      if (contentEl && !contentEl.contains(e.target)) {
        onDismiss();
      }
    },
    [contentEl, onDismiss]
  );

  useEventListener(document, 'mousedown', handleMousedown);
  useEventListener(contentEl, 'keydown', e => {
    switch (e.keyCode) {
      case KEYCODE_ESC: {
        e.stopPropagation();

        onDismiss();
        break;
      }
      case KEYCODE_TAB: {
        if (!contentEl) {
          return;
        }

        scopeTab(contentEl, e);
        break;
      }
      default:
      // noop
    }
  });

  // catch blur event of modal children and focus on the modal root
  // instead of focusing <body> so that keydown evenets added to the modal
  // root can be used.
  // useEventListener(document, 'focusout', e => {
  //   if (isOpen && contentEl) {
  //     if (contentEl.contains(e.target)) {
  //       e.preventDefault();
  //       contentEl.focus();
  //     }
  //   }
  // });

  useEffect(() => {
    const pageEl = document.body;

    if (isOpen) {
      scrollbarWidth = getScrollbarWidth();
      pageEl.classList.add('no-scroll');
      // Set padding-right if scrollbar was hidden
      pageEl.style.paddingRight = `${scrollbarWidth}px`;

      focusContent();
    } else {
      pageEl.classList.remove('no-scroll');
      pageEl.style.paddingRight = '';
    }
  }, [isOpen, focusContent]);

  if (!isOpen) {
    return null;
  }

  const modalRoot = document.getElementById('modal-root');

  return ReactDOM.createPortal(
    <div className={classnames(styles.overlay, overlayClassName)}>
      <div
        id={modalId}
        className={classnames(styles.content, modalClassName, {
          [styles.regular]: size === 'regular',
          [styles.small]: size === 'small'
        })}
        ref={el => {
          setContentEl(el);
        }}
        tabIndex={-1}
        role="dialog"
        aria-labelledby={ariaLabelledBy}
        aria-describedby={ariaDescribedBy}
        aria-modal="true"
      >
        {children}
      </div>
    </div>,
    modalRoot
  );
};

export default withRouter(Modal);

export { default as Header } from './Header';
export { default as Body } from './Body';
