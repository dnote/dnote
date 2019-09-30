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

/* eslint-disable react/button-has-type */

import React from 'react';
import classnames from 'classnames';

import SpinnerIcon from '../../Icons/Spinner';
import styles from './Button.scss';

type ButtonType = 'button' | 'submit' | 'reset';

interface Props {
  type: ButtonType;
  kind: string;
  size: string;
  children: React.ReactNode;
  id?: string;
  className?: string;
  isBusy?: boolean;
  stretch?: boolean;
  disabled?: boolean;
  onClick?: () => void;
  tabIndex?: number;
}

const Button: React.SFC<Props> = ({
  id,
  type,
  kind,
  size,
  children,
  className,
  isBusy,
  stretch,
  disabled,
  onClick,
  tabIndex
}) => {
  return (
    <button
      id={id}
      type={type}
      className={classnames(
        className,
        'button',
        `button-${kind}`,
        `button-${size}`,
        {
          [styles.busy]: isBusy,
          'button-stretch': stretch
        }
      )}
      disabled={isBusy || disabled}
      onClick={onClick}
      tabIndex={tabIndex}
    >
      <span className={styles.content}>{children}</span>

      {isBusy && (
        <SpinnerIcon
          width="16"
          height="16"
          className={styles.spinner}
          fill="white"
        />
      )}
    </button>
  );
};

Button.defaultProps = {
  size: 'normal'
};

export default Button;
