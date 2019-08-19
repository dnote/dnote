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
      <line
        y1="-1"
        x2="17.2395"
        y2="-1"
        transform="matrix(0.769979 0.638069 -0.653774 0.75669 3 11.5779)"
        stroke={fill}
        strokeWidth="4"
      />
      <line
        y1="-1"
        x2="17.2045"
        y2="-1"
        transform="matrix(0.739684 -0.672954 0.688196 0.725524 16.2741 22.5779)"
        stroke={fill}
        strokeWidth="4"
      />
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#a2a2a2',
  width: 32,
  height: 32
};

export default Icon;
