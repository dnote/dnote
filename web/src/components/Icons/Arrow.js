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
      width={w}
      height={h}
      viewBox="0 0 19 15"
      fill="none"
      className={className}
    >
      <path
        d="M1.27713 6.79289C0.886605 7.18342 0.886605 7.81658 1.27713 8.20711L7.64109 14.5711C8.03161 14.9616 8.66478 14.9616 9.0553 14.5711C9.44583 14.1805 9.44583 13.5474 9.0553 13.1569L3.39845 7.5L9.0553 1.84315C9.44583 1.45262 9.44583 0.819456 9.0553 0.428932C8.66478 0.0384079 8.03161 0.0384079 7.64109 0.428932L1.27713 6.79289ZM18.0522 6.5L1.98424 6.5V8.5L18.0522 8.5V6.5Z"
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
