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

interface Props extends IconProps {
  id: string;
}

const Icon = ({ fill, width, height, className, id }: Props) => {
  const h = `${height}px`;
  const w = `${width}px`;
  const clipPathId = `${id}-clip0`;

  return (
    <svg
      width={w}
      height={h}
      viewBox="0 0 32 32"
      fill="none"
      className={className}
    >
      <g clipPath={`url(#${clipPathId})`}>
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M19.875 28H4.125C2.67534 28 1.5 26.6568 1.5 25V7C1.5 5.34325 2.67534 4 4.125 4H19.875C21.3247 4 22.5 5.34325 22.5 7V11V15L19.3493 28C19.3493 29.6568 21.3247 28 19.875 28ZM6.75 7H4.65066C4.35928 7 4.125 7.26775 4.125 7.6V24.3993C4.125 24.7315 4.35928 25 4.65066 25H6.75V7ZM19.875 7.6C19.875 7.26775 19.6407 7 19.3493 7H17.25V14.1955C17.25 14.5247 17.0833 14.6042 16.8786 14.3702L14.9971 12.2192C14.7923 11.9852 14.459 11.9852 14.2536 12.2192L12.3721 14.3702C12.1667 14.6042 12 14.5247 12 14.1955V7H8.0625V25H19.3493C19.3493 31 19.3493 28 19.3493 28L22.5 15V11.3H19.875V7.6Z"
          fill={fill}
        />
      </g>
      <line
        x1="23.5"
        y1="13"
        x2="23.5"
        y2="28"
        stroke={fill}
        strokeWidth="2.5"
      />
      <line
        x1="16"
        y1="20.25"
        x2="31"
        y2="20.25"
        stroke={fill}
        strokeWidth="2.5"
      />
      <defs>
        <clipPath id={clipPathId}>
          <rect
            width="24"
            height="24"
            fill="white"
            transform="translate(0 4)"
          />
        </clipPath>
      </defs>
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000101',
  width: 32,
  height: 32
};

export default Icon;
