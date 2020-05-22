/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
import Helmet from 'react-helmet';
import { RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';

import { SettingSections } from 'web/libs/paths';
import Account from './Account';
import Sidebar from './Sidebar';
import About from './About';
import Notifications from './Notifications';
import styles from './Settings.scss';

function renderContent(section: string): React.ReactNode {
  if (section === SettingSections.account) {
    return <Account />;
  }
  if (section === SettingSections.about) {
    return <About />;
  }
  if (section === SettingSections.notifications) {
    return <Notifications />;
  }

  return <div>Not found</div>;
}

interface Match {
  section: string;
}

interface Props extends RouteComponentProps<Match> {}

const Settings: React.FunctionComponent<Props> = ({ match }) => {
  const { params } = match;
  const { section } = params;

  return (
    <div className="page page-mobile-full">
      <Helmet>
        <meta name="description" content="Dnote settings" />
      </Helmet>

      <div className="container mobile-fw">
        <div className={classnames('page-header', styles.header)}>
          <h1 className="page-heading">Settings</h1>
        </div>

        <div className="row">
          <div className="col-12 col-md-12 col-lg-3">
            <Sidebar />
          </div>

          <div className="col-12 col-md-12 col-lg-9">
            {renderContent(section)}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Settings;
