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
      viewBox="0 0 226 226"
      xmlSpace="preserve"
      width={w}
      height={h}
      fill={fill}
      className={className}
    >
      <path
        d="M65.307,83.651L44.782,22.916C63.74,8.537,87.371,0,113,0c42.059,0,78.748,22.98,98.204,57.067h-95.539
C114.781,57.026,113.894,57,113,57C92.836,57,75.167,67.662,65.307,83.651z M169,113c0,18.881-9.354,35.566-23.669,45.71
l-70.927,60.523C86.447,223.61,99.444,226,113,226c62.408,0,113-50.592,113-113c0-12.563-2.053-24.644-5.837-35.933h-64.221
C164.089,86.793,169,99.321,169,113z M54.658,209.79l48.776-41.621c-22.018-3.792-39.63-20.43-44.845-41.927L28.691,37.771
C10.85,57.75,0,84.106,0,113C0,154.061,21.902,190.003,54.658,209.79z M113,148c19.33,0,35-15.67,35-35s-15.67-35-35-35
s-35,15.67-35,35S93.67,148,113,148z"
      />
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000',
  width: 64,
  height: 64
};

export default Icon;
