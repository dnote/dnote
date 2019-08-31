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
MIT License

Copyright (c) 2019 GitHub Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
      viewBox="0 0 12 16"
      fill="none"
      className={className}
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M11 6H1C0.45 6 0 6.45 0 7V9C0 9.55 0.45 10 1 10H11C11.55 10 12 9.55 12 9V7C12 6.45 11.55 6 11 6ZM2 9H1V7H2V9ZM4 9H3V7H4V9ZM6 9H5V7H6V9ZM8 9H7V7H8V9ZM11 1H1C0.45 1 0 1.45 0 2V4C0 4.55 0.45 5 1 5H11C11.55 5 12 4.55 12 4V2C12 1.45 11.55 1 11 1ZM2 4H1V2H2V4ZM4 4H3V2H4V4ZM6 4H5V2H6V4ZM8 4H7V2H8V4ZM11 3H10V2H11V3ZM11 11H1C0.45 11 0 11.45 0 12V14C0 14.55 0.45 15 1 15H11C11.55 15 12 14.55 12 14V12C12 11.45 11.55 11 11 11ZM2 14H1V12H2V14ZM4 14H3V12H4V14ZM6 14H5V12H6V14ZM8 14H7V12H8V14Z"
        fill={fill}
      />
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000',
  width: 32,
  height: 32
};

export default Icon;
