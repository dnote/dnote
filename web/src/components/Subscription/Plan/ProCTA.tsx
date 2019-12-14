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
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import {
  getSubscriptionCheckoutPath,
  getSettingsPath,
  SettingSections
} from 'web/libs/paths';
import { UserData } from 'jslib/operations/types';

interface Props {
  user: UserData;
}

const ProCTA: React.FunctionComponent<Props> = ({ user }) => {
  if (user && user.pro) {
    return (
      <Link
        to={getSettingsPath(SettingSections.billing)}
        className="button button-large button-third-outline button-stretch"
      >
        Manage Your Plan
      </Link>
    );
  }

  return (
    <Link
      id="T-unlock-pro-btn"
      className={classnames('button button-large button-third button-stretch')}
      to={getSubscriptionCheckoutPath()}
    >
      Upgrade
    </Link>
  );
};

export default ProCTA;
