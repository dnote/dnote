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

import Plan from './internal';
import ServerIcon from '../../Icons/Server';
import GlobeIcon from '../../Icons/Globe';

const proPerks = [
  {
    id: 'hosted',
    icon: <ServerIcon width="16" height="16" fill="#4d4d8b" />,
    value: 'Fully hosted and managed'
  },
  {
    id: 'support',
    icon: <GlobeIcon width="16" height="16" fill="#4d4d8b" />,
    value: 'Support the Dnote community and development'
  }
];

function ProPlan({ wrapperClassName, ctaContent, bottomContent }) {
  return (
    <Plan
      name="Pro"
      price="$3"
      interval="month"
      perks={proPerks}
      wrapperClassName={wrapperClassName}
      ctaContent={ctaContent}
      bottomContent={bottomContent}
    />
  );
}

export default ProPlan;
