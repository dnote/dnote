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
      viewBox="0 0 32 32"
      width={w}
      height={h}
      fill={fill}
      className={className}
    >
      <g
        id="g854"
        transform="matrix(0.97739063,0,0,0.97739063,-7.0249696,14.045219)"
      >
        <g transform="translate(-0.2043139,-28)" id="g849">
          <path
            d="m 29.372796,28.05916 c 3.65179,7.479569 4.39975,14.519163 1.62791,15.927083 -2.85984,1.45191 -8.35952,-3.65179 -12.31929,-11.351347 -3.95977,-7.743554 -4.839722,-15.179126 -1.93589,-16.631042 1.18793,-0.615964 2.90383,0 4.83972,1.495914 M 11.333839,32.854883 C 9.2219596,31.754947 7.9900306,30.435023 8.0780256,29.027104 8.2100186,25.815289 15.381606,23.615415 24.049106,24.05539 c 8.6235,0.439975 15.5311,3.387805 15.39911,6.59962 -0.13199,1.319924 -1.49591,2.551853 -3.69579,3.431802 m -15.09113,7.523571 c -2.33186,1.6719 -4.39974,2.28786 -5.71967,1.45191 -2.639843,-1.7599 -0.96794,-9.063476 3.73979,-16.27906 4.70773,-7.215584 10.69138,-11.835318 13.37523,-10.163414 1.31992,0.879949 1.58391,2.859835 1.05594,5.499683"
            fill="none"
            stroke={fill}
            strokeWidth="1.36392128"
            strokeLinecap="round"
            className="path1"
          />
          <path
            d="m 25.765006,29.582024 c 0.21999,1.05594 -0.43997,2.111879 -1.49591,2.331866 -1.05594,0.219987 -2.06788,-0.439975 -2.28787,-1.495914 -0.21999,-1.055939 0.43997,-2.111878 1.49591,-2.331865 1.01194,-0.219988 2.06788,0.439974 2.28787,1.495913"
            fill={fill}
            strokeWidth="0.43997464"
            className="path2"
          />
        </g>
      </g>
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000',
  width: 131,
  height: 120
};

export default Icon;
