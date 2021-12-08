/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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

import styles from './index.scss';

type Position = 'top' | 'bottom';

interface Props {
  position?: Position;
  wrapperClassName?: string;
}

const PageToolbar: React.FunctionComponent<Props> = ({
  position,
  wrapperClassName,
  children
}) => {
  return (
    <div
      className={classnames(styles.wrapper, wrapperClassName, {
        [styles.bottom]: position === 'bottom'
      })}
    >
      {children}
    </div>
  );
};

export default PageToolbar;
