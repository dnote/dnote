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

const Icon = ({ fill, width, height, className }) => {
  const h = `${height}px`;
  const w = `${width}px`;

  return (
    <svg
      width={w}
      height={h}
      viewBox="0 0 16 16"
      fill="none"
      className={className}
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M1 4.27001V11.74C1 12.19 1.3 12.58 1.75 12.71L8.25 14.44C8.41 14.49 8.59 14.49 8.75 14.44L15.25 12.71C15.7 12.58 16 12.19 16 11.74V4.27001C16 3.82001 15.7 3.43001 15.25 3.30001L8.75 1.56001C8.59 1.53001 8.41 1.53001 8.25 1.56001L1.75 3.30001C1.3 3.43001 1 3.82001 1 4.27001ZM8 13.36L2 11.77V5.00001L8 6.61001V13.36ZM2 4.00001L4.5 3.33001L11 5.06001L8.5 5.73001L2 4.00001ZM15 11.77L9 13.36V6.61001L11 6.06001V8.50001L13 7.97001V5.53001L15 5.00001V11.77ZM13 4.53L6.5 2.8L8.5 2.27L15 4L13 4.53Z"
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
