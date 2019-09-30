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
      viewBox="0 0 14 16"
      fill="none"
      className={className}
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M7 1C3.14 1 0 4.14 0 8C0 11.86 3.14 15 7 15C7.48 15 7.94 14.95 8.38 14.86C8.21 14.78 8.18 14.13 8.36 13.77C8.55 13.36 9.17 12.32 8.56 11.97C7.95 11.62 8.12 11.47 7.75 11.06C7.38 10.65 7.53 10.59 7.5 10.48C7.42 10.14 7.86 9.59 7.89 9.54C7.91 9.48 7.91 9.27 7.89 9.21C7.89 9.13 7.62 8.99 7.55 8.98C7.49 8.98 7.44 9.09 7.35 9.11C7.26 9.13 6.85 8.86 6.76 8.78C6.67 8.7 6.62 8.55 6.49 8.44C6.36 8.31 6.35 8.41 6.16 8.33C5.97 8.25 5.36 8.02 4.88 7.85C4.4 7.66 4.36 7.38 4.36 7.19C4.34 6.99 4.06 6.72 3.94 6.52C3.8 6.32 3.78 6.05 3.74 6.11C3.7 6.17 3.99 6.89 3.94 6.92C3.89 6.94 3.78 6.72 3.64 6.54C3.5 6.35 3.78 6.45 3.34 5.59C2.9 4.73 3.48 4.29 3.51 3.84C3.54 3.39 3.89 4.01 3.7 3.71C3.51 3.41 3.7 2.82 3.56 2.6C3.43 2.38 2.68 2.85 2.68 2.85C2.7 2.63 3.37 2.27 3.84 1.93C4.31 1.59 4.62 1.87 5 1.98C5.39 2.11 5.41 2.07 5.28 1.93C5.15 1.8 5.34 1.76 5.64 1.8C5.92 1.85 6.02 2.21 6.47 2.16C6.94 2.13 6.52 2.25 6.58 2.38C6.64 2.51 6.52 2.49 6.2 2.68C5.9 2.88 6.22 2.9 6.75 3.29C7.28 3.68 7.13 3.04 7.06 2.74C6.99 2.44 7.45 2.68 7.45 2.68C7.78 2.9 7.72 2.7 7.95 2.76C8.18 2.82 8.86 3.4 8.86 3.4C8.03 3.84 8.55 3.88 8.69 3.99C8.83 4.1 8.41 4.29 8.41 4.29C8.24 4.12 8.22 4.31 8.11 4.37C8 4.43 8.09 4.59 8.09 4.59C7.53 4.68 7.65 5.28 7.67 5.42C7.67 5.56 7.29 5.78 7.2 6C7.11 6.2 7.45 6.64 7.26 6.66C7.07 6.69 6.92 6 5.95 6.25C5.65 6.33 5.01 6.66 5.36 7.33C5.72 8.02 6.28 7.14 6.47 7.24C6.66 7.34 6.41 7.77 6.45 7.79C6.49 7.81 6.98 7.81 7.01 8.4C7.04 8.99 7.78 8.93 7.93 8.95C8.1 8.95 8.63 8.51 8.7 8.5C8.76 8.47 9.08 8.22 9.73 8.59C10.39 8.95 10.71 8.9 10.93 9.06C11.15 9.22 11.01 9.53 11.21 9.64C11.41 9.75 12.27 9.61 12.49 9.95C12.71 10.29 11.61 12.04 11.27 12.23C10.93 12.42 10.79 12.87 10.43 13.15C10.07 13.43 9.62 13.79 9.16 14.06C8.75 14.29 8.69 14.72 8.5 14.86C11.64 14.16 13.98 11.36 13.98 8.02C13.98 4.16 10.84 1.02 6.98 1.02L7 1ZM8.64 7.56C8.55 7.59 8.36 7.78 7.86 7.48C7.38 7.18 7.05 7.25 7 7.2C7 7.2 6.95 7.09 7.17 7.06C7.61 7.01 8.15 7.47 8.28 7.47C8.41 7.47 8.47 7.34 8.69 7.42C8.91 7.5 8.74 7.55 8.64 7.56ZM6.34 1.7C6.29 1.67 6.37 1.62 6.43 1.56C6.46 1.53 6.45 1.45 6.48 1.42C6.59 1.31 7.09 1.17 7 1.45C6.89 1.72 6.42 1.75 6.34 1.7ZM7.57001 2.59C7.38001 2.57 6.99001 2.54 7.05001 2.45C7.35001 2.17 6.96001 2.07 6.71001 2.07C6.46001 2.05 6.37001 1.91 6.49001 1.88C6.61001 1.85 7.10001 1.9 7.19001 1.96C7.27001 2.02 7.71001 2.21 7.74001 2.34C7.76001 2.47 7.74001 2.59 7.57001 2.59ZM9.03999 2.54C8.89999 2.63 8.20999 2.13 8.08999 2.02C7.52999 1.54 7.19999 1.71 7.08999 1.61C6.97999 1.51 7.00999 1.42 7.19999 1.27C7.38999 1.12 7.88999 1.33 8.19999 1.36C8.49999 1.39 8.85999 1.63 8.85999 1.91C8.87999 2.16 9.18999 2.41 9.04999 2.54H9.03999Z"
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
