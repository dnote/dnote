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
import { withRouter, RouteComponentProps } from 'react-router-dom';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';

import { getPathFromLocation } from 'jslib//helpers/url';
import { AppState, RemoteData } from '../store';
import { UserData } from '../store/auth';
import { useSelector } from '../store';

// userOnly returns a HOC that redirects to Login page if user is not logged in
export default function(
  Component: React.ComponentType,
  guestPath: string = '/login'
) {
  interface Props extends RouteComponentProps {}

  const HOC: React.SFC<Props> = props => {
    const { location } = props;

    const { userData } = useSelector(state => {
      return {
        userData: state.auth.user
      };
    });

    const isGuest = userData.isFetched && !userData.data.uuid;
    if (isGuest) {
      const referrer = getPathFromLocation(location);

      const dest = `${guestPath}?referrer=${encodeURIComponent(referrer)}`;

      return <Redirect to={dest} />;
    }

    return <Component {...props} />;
  };

  return withRouter(HOC);
}
