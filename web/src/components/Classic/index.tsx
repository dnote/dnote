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
import { Switch, Route } from 'react-router';
import { Redirect } from 'react-router-dom';

import { useDispatch, useSelector } from '../../store';
import ClassicLogin from './Login';
import { setMessage } from '../../store/ui';
import ClassicSetPassword from './SetPassword';
import ClassicDecrypt from './Decrypt';
import {
  ClassicMigrationSteps,
  getClassicMigrationPath,
  getHomePath,
  homePathDef
} from 'web/libs/paths';

interface Props {}

const Classic: React.SFC<Props> = () => {
  const { user } = useSelector(state => {
    return {
      user: state.auth.user
    };
  });
  const dispatch = useDispatch();

  if (!user.isFetched) {
    return <div>Loading</div>;
  }

  const userData = user.data;
  const loggedIn = userData.uuid !== '';

  if (loggedIn && !userData.classic) {
    dispatch(
      setMessage({
        message:
          'You are already using the latest Dnote and do not have to migrate.',
        kind: 'info',
        path: homePathDef
      })
    );

    return <Redirect to={getHomePath()} />;
  }

  return (
    <div className="container">
      <Switch>
        <Route
          path={getClassicMigrationPath(ClassicMigrationSteps.login)}
          exact
          component={ClassicLogin}
        />
        <Route
          path={getClassicMigrationPath(ClassicMigrationSteps.setPassword)}
          exact
          component={ClassicSetPassword}
        />
        <Route
          path={getClassicMigrationPath(ClassicMigrationSteps.decrypt)}
          exact
          component={ClassicDecrypt}
        />
      </Switch>
    </div>
  );
};

export default Classic;
