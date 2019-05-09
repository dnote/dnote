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
    <svg height={h} width={w} viewBox="0 0 24 24" className={className}>
      <g fill="none" fillRule="evenodd" id="miu" stroke="none" strokeWidth="1">
        <g transform="translate(-827.000000, -371.000000)">
          <g transform="translate(215.000000, 119.000000)" />
          <path
            d="M828,373.002462 L828,391.997538 C828,392.551183 828.450975,393 828.990778,393 L832.009222,393 C832.556414,393 833,392.560542 833,391.997538 L833,373.002462 C833,372.448817 832.549025,372 832.009222,372 L828.990778,372 C828.443586,372 828,372.439458 828,373.002462 Z M834,373.002462 L834,391.997538 C834,392.551183 834.450975,393 834.990778,393 L838.009222,393 C838.556414,393 839,392.560542 839,391.997538 L839,373.002462 C839,372.448817 838.549025,372 838.009222,372 L834.990778,372 C834.443586,372 834,372.439458 834,373.002462 Z M839.627042,373.97313 L844.543329,392.320965 C844.686623,392.855744 845.238394,393.172548 845.759803,393.032837 L848.675397,392.251606 C849.203943,392.109982 849.518674,391.570689 849.372958,391.02687 L844.456671,372.679035 C844.313377,372.144256 843.761606,371.827452 843.240197,371.967163 L840.324603,372.748394 C839.796057,372.890018 839.481326,373.429311 839.627042,373.97313 Z"
            fill={fill}
          />
        </g>
      </g>
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000',
  width: 32,
  height: 32
};

export default Icon;
