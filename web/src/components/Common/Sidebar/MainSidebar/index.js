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

import MainSidebar from './MainSidebar';
import { mainSidebarPaths } from 'web/libs/paths';

export default ({ isEditor, demo }) => {
  return (
    <Switch>
      <Route
        exact
        path={['/notes/:noteUUID', '/demo/notes/:noteUUID']}
        render={() => {
          if (isEditor) {
            return <MainSidebar demo={demo} />;
          }

          return null;
        }}
      />

      <Route
        exact
        path={mainSidebarPaths}
        render={() => {
          return <MainSidebar demo={demo} />;
        }}
      />
    </Switch>
  );
};
