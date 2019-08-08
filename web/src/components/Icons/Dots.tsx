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
      height={h}
      viewBox="0 0 128 128"
      width={w}
      fill={fill}
      className={className}
    >
      <circle cx="28.75" cy="64" r="12" />
      <circle cx="68.75" cy="64" r="12" />
      <circle cx="108.75" cy="64" r="12" />
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#2F3435',
  width: 20,
  height: 4
};

export default Icon;
