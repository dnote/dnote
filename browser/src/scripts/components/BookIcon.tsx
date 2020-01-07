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

import React from 'react';

const Icon = ({ fill, width, height, className }) => {
  const h = `${height}px`;
  const w = `${width}px`;

  return (
    <svg
      viewBox="0 0 32 32"
      height={h}
      width={w}
      fill={fill}
      className={className}
    >
      <g transform="translate(240 0)">
        <path d="M-211,4v26h-24c-1.104,0-2-0.895-2-2s0.896-2,2-2h22V0h-22c-2.209,0-4,1.791-4,4v24c0,2.209,1.791,4,4,4h26V4H-211z    M-235,8V2h20v22h-20V8z M-219,6h-12V4h12V6z M-223,10h-8V8h8V10z M-227,14h-4v-2h4V14z" />
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
