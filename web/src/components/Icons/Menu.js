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
    <svg width={h} height={w} viewBox="0 0 50 50" className={className}>
      <g fill={fill} stroke="none" strokeWidth="1" fillRule="evenodd">
        <rect x="6" y="9" width="37.9969741" height="6" />
        <rect x="6" y="22" width="28" height="6" />
        <rect x="6" y="35" width="37.9969741" height="6" />
      </g>
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#686868',
  width: 50,
  height: 50
};

export default Icon;
