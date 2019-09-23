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

import React, { Fragment } from 'react';
import Helmet from 'react-helmet';
import { RouteComponentProps } from 'react-router-dom';

import Account from './Account';
import Sidebar from './Sidebar';
import { SettingSections } from 'web/libs/paths';
import Notification from './Notification';
import Billing from './Billing';

function renderContent(section: string): React.ReactNode {
  if (section === SettingSections.account) {
    return <Account />;
  }
  if (section === SettingSections.notification) {
    return <Notification />;
  }
  if (section === SettingSections.billing) {
    return <Billing />;
  }

  return <div>Not found</div>;
}

interface Match {
  section: string;
}

interface Props extends RouteComponentProps<Match> {}

const Settings: React.SFC<Props> = ({ match }) => {
  const { params } = match;
  const { section } = params;

  return (
    <Fragment>
      <Helmet>
        <meta name="description" content="Dnote settings" />
      </Helmet>

      <div className="container">
        <div className="row">
          <div className="col-12 col-md-12 col-lg-3">
            <Sidebar />
          </div>

          <div className="col-12 col-md-12 col-lg-9">
            {renderContent(section)}
          </div>
        </div>
      </div>
    </Fragment>
  );
};

export default Settings;
