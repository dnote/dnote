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

interface Props extends IconProps {
  id: string;
}

const Icon = ({ fill, width, height, id }: Props) => {
  const h = `${height}px`;
  const w = `${width}px`;
  const clipPathId = `${id}-clip0`;

  return (
    <svg width={w} height={h} viewBox="0 0 32 32" fill="none">
      <g clipPath={`url(#${clipPathId})`}>
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M30 30.002C30 31.106 29.105 32.002 28 32.002H4C2.895 32.002 2 31.106 2 30.002V2.002C2 0.898 2.895 0.002 4 0.002L21.158 0C21.599 0 22.213 0.255 22.523 0.566L29.435 7.476C29.746 7.788 30.001 8.402 30.001 8.842L30 30.002ZM26 10.002H22C20.895 10.002 20 9.106 20 8.002V4.002H6.801C6.357 4.002 6 4.359 6 4.802V27.201C6 27.644 6.357 28.002 6.801 28.002H25.199C25.643 28.002 26 27.644 26 27.201V10.002ZM22 24.002H10C8.895 24.002 8 23.106 8 22.002C8 20.897 8.895 20.002 10 20.002H22C23.105 20.002 24 20.897 24 22.002C24 23.106 23.105 24.002 22 24.002ZM22 18.002H10C8.895 18.002 8 17.106 8 16.002C8 14.898 8.895 14.002 10 14.002H22C23.105 14.002 24 14.898 24 16.002C24 17.106 23.105 18.002 22 18.002Z"
          fill={fill}
        />
      </g>
      <defs>
        <clipPath id={clipPathId}>
          <rect width="32" height="32" fill="white" />
        </clipPath>
      </defs>
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000101',
  width: 32,
  height: 32
};

export default Icon;
