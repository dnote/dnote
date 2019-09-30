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

import { IconProps } from './types';

const Icon = ({ fill, width, height, className }: IconProps) => {
  const h = `${height}px`;
  const w = `${width}px`;

  return (
    <svg viewBox="0 0 20 20" height={h} width={w} className={className}>
      <g fill="none" fillRule="evenodd" stroke="none" strokeWidth="1">
        <g fill={fill} transform="translate(-86.000000, -2.000000)">
          <g id="account-circle" transform="translate(86.000000, 2.000000)">
            <path d="M10,0 C4.5,0 0,4.5 0,10 C0,15.5 4.5,20 10,20 C15.5,20 20,15.5 20,10 C20,4.5 15.5,0 10,0 L10,0 Z M10,3 C11.7,3 13,4.3 13,6 C13,7.7 11.7,9 10,9 C8.3,9 7,7.7 7,6 C7,4.3 8.3,3 10,3 L10,3 Z M10,17.2 C7.5,17.2 5.3,15.9 4,14 C4,12 8,10.9 10,10.9 C12,10.9 16,12 16,14 C14.7,15.9 12.5,17.2 10,17.2 L10,17.2 Z" />
          </g>
        </g>
      </g>
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000',
  width: 32,
  height: 32
};

export default Icon;
