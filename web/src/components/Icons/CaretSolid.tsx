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
      enableBackground="new 0 0 29 14"
      height={h}
      viewBox="0 0 29 14"
      width={w}
      className={className}
    >
      <polygon fill={fill} points="0.15,0 14.5,14.35 28.85,0 " />
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#231F20',
  width: 12,
  height: 8
};

export default Icon;
