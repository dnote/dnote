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

const Icon = ({ className, width, height }) => {
  return (
    <svg
      version="1.1"
      x="0px"
      y="0px"
      viewBox="0 0 512 512"
      width={width}
      height={height}
      className={className}
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
    </svg>
  );
};

Icon.defaultProps = {
  width: 16,
  height: 16
};

export default Icon;
