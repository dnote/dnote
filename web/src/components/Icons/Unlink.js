/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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

const Icon = ({ className, width, height }) => {
  return (
    <svg
      version="1.1"
      x="0px"
      y="0px"
      viewBox="0 0 512 512"
      className={className}
      width={width}
      height={height}
    >
      <g>
        <g>
          <path
            d="M472.503,39.494c-52.659-52.659-138.33-52.659-190.99,0l-78.469,78.469c-13.693,13.693-24.788,31.163-31.939,52.492
			c-0.12,0.04-0.236,0.101-0.356,0.141c-16.264,48.927-3.929,102.132,32.295,138.356c10.823,10.823,23.557,19.735,37.349,26.314
			l5.093-5.093c17.614-17.614,17.614-46.049,0-63.664c-29.285-29.285-29.285-76.821,0-106.105l78.469-78.469
			c29.285-29.285,76.82-29.285,106.105,0c29.285,29.285,29.285,76.82,0,106.105l-59.582,59.582
			c7.945,24.565,9.628,50.578,5.783,75.897c6.275-4.374,12.313-9.108,17.772-14.567l78.469-78.469
			C525.162,177.824,525.162,92.152,472.503,39.494z"
          />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M309.151,202.847c-10.823-10.823-23.555-19.736-37.56-26.103l-4.882,4.882c-17.613,17.613-17.614,46.049,0,63.664
			c29.285,29.285,29.285,76.821,0,106.105c-24.788,24.788-55.03,55.035-78.664,78.669c-29.285,29.285-76.82,29.285-106.105,0
			c-29.285-29.285-29.285-76.82,0-106.105l59.777-59.782c-7.945-24.565-9.627-50.579-5.782-75.898
			c-6.276,4.376-12.314,9.109-17.774,14.568l-78.664,78.669c-52.659,52.659-52.659,138.33,0,190.99
			c52.659,52.66,138.33,52.659,190.99,0l78.664-78.669c13.693-13.694,24.788-31.163,31.939-52.492
			c0.12-0.04,0.236-0.101,0.356-0.141C357.71,292.276,345.375,239.071,309.151,202.847z"
          />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M436.477,415.057l-21.221-21.221c-5.865-5.865-15.357-5.865-21.221,0c-5.865,5.865-5.865,15.357,0,21.221l21.221,21.221
			c5.865,5.865,15.357,5.865,21.221,0C442.341,430.413,442.341,420.922,436.477,415.057z"
          />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M118.161,96.741L96.939,75.52c-5.865-5.865-15.357-5.865-21.221,0c-5.865,5.865-5.865,15.357,0,21.221l21.221,21.221
			c5.865,5.865,15.357,5.865,21.222,0S124.026,102.605,118.161,96.741z"
          />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M491.892,303.542l-28.993-7.772c-8.01-2.145-16.237,2.601-18.382,10.611c-2.146,8.059,2.655,16.253,10.611,18.383
			l28.993,7.771c8.01,2.145,16.237-2.601,18.382-10.611C504.626,313.914,499.901,305.646,491.892,303.542z"
          />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M57.067,187.035l-28.993-7.772c-8.01-2.145-16.237,2.601-18.382,10.611c-2.147,8.057,2.654,16.252,10.611,18.382
			l28.993,7.772c8.01,2.145,16.237-2.601,18.382-10.611C69.812,197.397,65.046,189.127,57.067,187.035z"
          />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M332.733,483.922l-7.771-28.993c-2.155-8.041-10.372-12.756-18.382-10.611c-8.01,2.145-12.756,10.372-10.611,18.382
			l7.772,28.993c2.13,7.957,10.325,12.757,18.382,10.611C330.133,500.159,334.878,491.932,332.733,483.922z"
          />
        </g>
      </g>
      <g>
        <g>
          <path
            d="M216.225,49.097l-7.772-28.993c-2.135-8.041-10.372-12.735-18.382-10.611c-8.01,2.145-12.756,10.372-10.611,18.382
			l7.772,28.993c0.704,2.632,2.073,4.912,3.875,6.715c3.668,3.668,9.139,5.326,14.506,3.896
			C213.624,65.335,218.37,57.107,216.225,49.097z"
          />
        </g>
      </g>
    </svg>
  );
};

Icon.defaultProps = {
  width: 16,
  height: 16
};

export default Icon;
