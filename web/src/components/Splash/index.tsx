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

import Logo from '../Icons/Logo';
import styles from './Splash.module.scss';

function Splash() {
  return (
    <div className={styles.wrapper}>
      <Logo width={60} height={60} fill="#252833" />

      <svg
        version="1.0"
        width="40px"
        height="25px"
        viewBox="0 0 128 32"
        style={{ marginTop: '12px' }}
      >
        <circle
          fill="#848484"
          fillOpacity="1"
          cx="0"
          cy="0"
          r="12"
          transform="translate(16 16)"
        >
          <animateTransform
            attributeName="transform"
            type="scale"
            additive="sum"
            values="1;1.42;1;1;1;1;1;1;1;1"
            dur="1350ms"
            repeatCount="indefinite"
          />
        </circle>
        <circle
          fill="#848484"
          fillOpacity="1"
          cx="0"
          cy="0"
          r="12"
          transform="translate(64 16)"
        >
          <animateTransform
            attributeName="transform"
            type="scale"
            additive="sum"
            values="1;1;1;1;1.42;1;1;1;1;1"
            dur="1350ms"
            repeatCount="indefinite"
          />
        </circle>
        <circle
          fill="#848484"
          fillOpacity="1"
          cx="0"
          cy="0"
          r="12"
          transform="translate(112 16)"
        >
          <animateTransform
            attributeName="transform"
            type="scale"
            additive="sum"
            values="1;1;1;1;1;1;1;1.42;1;1"
            dur="1350ms"
            repeatCount="indefinite"
          />
        </circle>
      </svg>
    </div>
  );
}

export default Splash;
