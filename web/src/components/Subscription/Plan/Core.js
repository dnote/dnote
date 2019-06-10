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

import BoxIcon from '../../Icons/Box';
import Plan from './internal';

const selfHostedPerks = [
  {
    id: 'own-machine',
    icon: <BoxIcon width="16" height="16" fill="#6e6e6e" />,
    value: 'Host on your own machine'
  }
];

function Core({ wrapperClassName, ctaContent, bottomContent }) {
  return (
    <Plan
      name="Core"
      price="Free"
      perks={selfHostedPerks}
      wrapperClassName={wrapperClassName}
      ctaContent={ctaContent}
      bottomContent={bottomContent}
    />
  );
}

export default Core;
