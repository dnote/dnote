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

const Icon = ({ fill, width, height, className }) => {
  const h = `${height}px`;
  const w = `${width}px`;

  return (
    <svg
      enableBackground="new 0 0 32 32"
      height={h}
      width={w}
      version="1.1"
      viewBox="0 0 32 32"
      className={className}
    >
      <path
        d="M30,10l-2,18c0,2.209-1.791,4-4,4H8c-2.209,0-4-1.791-4-4L2,10  c-1.105,0-2-0.896-2-2c0-1.105,0.895-2,2-2h4h4V4c0-2.209,1.791-4,4-4h4c2.209,0,4,1.791,4,4v2h4h4c1.105,0,2,0.895,2,2  C32,9.104,31.105,10,30,10z M18,5.199V5V4.8C18,4.357,17.643,4,17.199,4h-2.398C14.357,4,14,4.357,14,4.8V5v0.199V6h0.801h2.398H18  V5.199z M25.199,10H6.801C6.357,10,6,10.357,6,10.8l2,16.399C8,27.641,8.357,28,8.801,28h14.398C23.643,28,24,27.641,24,27.199  L26,10.8C26,10.357,25.643,10,25.199,10z M20,24c-1.105,0-2-0.896-2-2v-6c0-1.104,0.895-2,2-2s2,0.896,2,2v6  C22,23.104,21.105,24,20,24z M12,24c-1.105,0-2-0.896-2-2v-6c0-1.104,0.895-2,2-2c1.104,0,2,0.896,2,2v6C14,23.104,13.104,24,12,24z  "
        fill={fill}
        fillRule="evenodd"
        clipRule="evenodd"
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
