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

const Icon = ({ width, height }) => {
  const h = `${height}px`;
  const w = `${width}px`;

  return (
    <svg
      enableBackground="new 0 0 500 500"
      height={h}
      viewBox="0 0 500 500"
      width={w}
    >
      <path d="M418.327,188.854H393.73v-13.627c0-80.443-65.183-145.895-145.274-145.895c-80.078,0-145.26,65.452-145.26,145.895v13.627  H81.581c-12.278,0-22.195,10.025-22.195,22.411v236.383c0,12.359,9.917,22.371,22.195,22.371h336.747  c12.278,0,22.195-10.012,22.195-22.371V211.266C440.522,198.88,430.605,188.854,418.327,188.854z M288.26,391.411h-76.611  l19.848-48.815c-12.737-6.639-21.453-19.834-21.453-35.188c0-22.047,17.837-39.932,39.911-39.932  c22.02,0,39.911,17.885,39.911,39.932c0,15.354-8.771,28.55-21.453,35.188L288.26,391.411z M342.338,188.854H154.589v-13.627  c0-52.095,42.124-94.502,93.867-94.502c51.758,0,93.882,42.407,93.882,94.502V188.854z" />
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000',
  width: 32,
  height: 32
};

export default Icon;
