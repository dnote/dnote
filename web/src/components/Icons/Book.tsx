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
    <svg
      width={w}
      height={h}
      viewBox="0 0 32 32"
      fill="none"
      className={className}
    >
      <title>Book</title>
      <desc>Icon depicting a book</desc>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M26.5 32H5.5C3.56712 32 2 30.209 2 28V4C2 1.791 3.56712 0 5.5 0H26.5C28.4329 0 30 1.791 30 4V28C30 30.209 28.4329 32 26.5 32ZM9 4H6.20088C5.81238 4 5.5 4.357 5.5 4.8V27.199C5.5 27.642 5.81238 28 6.20088 28H9V4ZM26.5 4.8C26.5 4.357 26.1876 4 25.7991 4H23V13.594C23 14.033 22.7778 14.139 22.5048 13.827L19.9961 10.959C19.7231 10.647 19.2786 10.647 19.0048 10.959L16.4961 13.827C16.2222 14.139 16 14.033 16 13.594V4H10.75V28H25.7991C26.1876 28 26.5 27.642 26.5 27.199V4.8Z"
        fill={fill}
      />
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000101',
  width: 32,
  height: 32
};

export default Icon;
