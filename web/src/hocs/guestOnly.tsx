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
import { connect } from 'react-redux';
import { Redirect, RouteComponentProps } from 'react-router-dom';
import hoistNonReactStatics from 'hoist-non-react-statics';

import { getReferrer } from 'jslib//helpers/url';
import { AppState, RemoteData } from '../store';
import { UserData } from '../store/auth';

function renderFallback(referrer?: string) {
  let destination;
  if (referrer) {
    destination = referrer;
  } else {
    destination = '/';
  }

  return <Redirect to={{ pathname: destination }} />;
}

// guestOnly returns a HOC that renders the given component only if user is not
// logged in
export default function(Component: React.ComponentType): React.ComponentType {
  interface Props extends RouteComponentProps {
    userData: RemoteData<UserData>;
  }

  const HOC: React.SFC<Props> = props => {
    const { userData, location } = props;

    const loggedIn = userData.isFetched && Boolean(userData.data.uuid);

    if (loggedIn) {
      const referrer = getReferrer(location);
      return renderFallback(referrer);
    }

    return <Component {...props} />;
  };

  // Copy over static methods
  hoistNonReactStatics(HOC, Component);

  function mapStateToProps(state: AppState) {
    return {
      userData: state.auth.user
    };
  }

  return connect(mapStateToProps)(HOC);
}
