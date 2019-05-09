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

import styles from './Splash.module.scss';

function Splash() {
  return (
    <div className={styles.wrapper}>
      <svg width="60" height="60" viewBox="0 0 32 30">
        <g>
          <path
            d="M27.7993 7.9354C23.7727 7.11689 20.126 5.10464 17.2593 2.11388L16.8326 1.66467C16.706 1.53538 16.526 1.46154 16.3393 1.46154H16.3326C16.146 1.46154 15.966 1.52922 15.8393 1.65851L15.186 2.32309C12.4593 5.08616 9.0393 6.99384 5.286 7.84308L4.84598 7.94156C4.55932 8.00308 4.34598 8.24309 4.33931 8.51998C4.09932 17.1538 9.61934 23.3754 12.006 25.6462C13.506 27.08 15.426 28.5323 16.326 28.5385C16.3326 28.5385 16.3326 28.5385 16.3326 28.5385C16.3393 28.5385 16.346 28.5385 16.346 28.5385C17.3593 28.5138 19.3393 26.9323 20.726 25.5846C23.0993 23.2769 28.586 16.9938 28.326 8.51998C28.3193 8.23693 28.0993 7.99692 27.7993 7.9354ZM17.1993 21.6647C17.086 21.7323 16.9593 21.7692 16.8326 21.7692C16.726 21.7692 16.6193 21.7446 16.5193 21.6954C16.2993 21.5846 16.166 21.3754 16.166 21.1538V8.84616C16.166 8.63695 16.2793 8.44006 16.4726 8.32925C16.6593 8.21229 16.8993 8.19997 17.106 8.28005L23.7727 11.0493C24.0393 11.16 24.1927 11.4123 24.166 11.677C24.1326 11.9292 23.386 17.9231 17.1993 21.6647Z"
            fill="#252833"
          />
        </g>
      </svg>
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
