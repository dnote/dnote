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

import React from 'react';
import classnames from 'classnames';

import CloseIcon from '../../Icons/Close';

import styles from './Flash.scss';

const TYPE_SUCCESS = 'success';
const TYPE_INFO = 'info';
const TYPE_WARNING = 'warning';
const TYPE_DANGER = 'danger';

const validTypes = [TYPE_SUCCESS, TYPE_INFO, TYPE_WARNING, TYPE_DANGER];

function validateType(kind) {
  return validTypes.indexOf(kind) > -1;
}

interface Props {
  id?: string;
  kind: string;
  when?: boolean;
  hasBorder?: boolean;
  noMargin?: boolean;
  onDismiss?: () => void;
  wrapperClassName?: string;
  contentClassName?: string;
  children: React.ReactNode;
}

const Flash: React.SFC<Props> = ({
  id,
  when,
  kind,
  children,
  hasBorder,
  onDismiss,
  noMargin,
  wrapperClassName,
  contentClassName
}) => {
  // If `when` prop is provided and is explicitly false, do not render
  if (when === false) {
    return null;
  }

  if (!validateType(kind)) {
    console.log(`Invalid kind ${kind}`);
  }

  const dismissable = Boolean(onDismiss);

  return (
    <div
      id={id}
      className={classnames(styles.wrapper, wrapperClassName, {
        [styles.success]: kind === TYPE_SUCCESS,
        [styles.info]: kind === TYPE_INFO,
        [styles.warning]: kind === TYPE_WARNING,
        [styles.danger]: kind === TYPE_DANGER,
        [styles.dismissable]: dismissable,
        [styles.border]: hasBorder,
        [styles.nomargin]: noMargin
      })}
    >
      <div className={contentClassName}>{children}</div>

      {dismissable && (
        <button
          type="button"
          className={classnames('button-no-ui', styles.dismiss)}
          aria-label="Dismiss this message"
          onClick={onDismiss}
        >
          <CloseIcon width={16} height={16} />
        </button>
      )}
    </div>
  );
};

Flash.defaultProps = {
  hasBorder: true
};

export default Flash;
