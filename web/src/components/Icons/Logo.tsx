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
      xmlns="http://www.w3.org/2000/svg"
      className={className}
    >
      <path
        d="M8.5968 1.5412C8.53682 1.54101 8.47683 1.54083 8.4167 1.57058C8.38671 1.57048 8.35657 1.60032 8.32658 1.60023C8.29659 1.60014 8.26645 1.62998 8.23645 1.62989C8.20646 1.62979 8.17632 1.65963 8.14618 1.68948C8.11619 1.68938 8.08605 1.71922 8.08605 1.71922C8.02592 1.74897 7.99578 1.77882 7.9355 1.8385L1.90756 7.8069C1.87742 7.83674 1.81714 7.89643 1.78685 7.95621C1.78671 7.98614 1.75671 7.98605 1.75657 8.01598C1.72643 8.04583 1.72628 8.07576 1.69614 8.1056C1.69599 8.13554 1.66585 8.16538 1.66571 8.19532C1.66556 8.22525 1.63542 8.25509 1.63527 8.28503C1.63498 8.3449 1.60469 8.40468 1.6044 8.46455L1.50125 29.4195C1.49875 29.9284 1.88673 30.3187 2.3966 30.3203L23.3912 30.3858C23.4512 30.386 23.5112 30.3862 23.5713 30.3565C23.6013 30.3566 23.6315 30.3267 23.6615 30.3268C23.6915 30.3269 23.7216 30.2971 23.7516 30.2972C23.7816 30.2972 23.8117 30.2674 23.8419 30.2376C23.8719 30.2377 23.902 30.2078 23.902 30.2078C23.9621 30.1781 23.9923 30.1482 24.0525 30.0885L30.0805 24.1201C30.2613 23.9411 30.3523 23.7318 30.3535 23.4923L30.4568 2.50748C30.4893 1.99867 30.1013 1.60829 29.5914 1.6067L8.5968 1.5412ZM8.58619 3.69656L8.56733 7.52832L23.2036 7.57398L25.9463 4.85836C26.308 4.50026 26.8778 4.50203 27.206 4.86229C27.5641 5.22264 27.5613 5.79142 27.1998 6.11958L24.4571 8.83521L24.3852 23.4438L28.2242 23.4558L23.1005 28.5289L3.36549 28.4673L3.46245 8.76971L8.58619 3.69656ZM8.55849 9.32445L8.48924 23.3942L22.5856 23.4382L22.6549 9.36843L8.55849 9.32445Z"
        fill={fill}
      />
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#2a2a2a',
  width: 32,
  height: 32
};

export default Icon;
