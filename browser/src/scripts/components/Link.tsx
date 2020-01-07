/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

/* eslint-disable react/jsx-props-no-spreading */

import React from 'react';

import { useDispatch } from '../store/hooks';
import { navigate } from '../store/location/actions';

interface Props {
  to: string;
  className: string;
  tabIndex?: number;
  onClick?: () => void;
}

const Link: React.FunctionComponent<Props> = ({
  to,
  children,
  onClick,
  ...restProps
}) => {
  const dispatch = useDispatch();

  return (
    <a
      href={`${to}`}
      onClick={e => {
        e.preventDefault();

        dispatch(navigate(to));

        if (onClick) {
          onClick();
        }
      }}
      {...restProps}
    >
      {children}
    </a>
  );
};

export default Link;
