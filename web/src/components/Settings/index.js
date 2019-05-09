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
import Helmet from 'react-helmet';

import Account from './Account';
import Notification from './Notification';
import Billing from './Billing';

import './module.scss';

const SectionAccount = 'account';
const SectionEmail = 'notification';
const SectionBilling = 'billing';

function renderContent(section) {
  if (section === SectionAccount) {
    return <Account />;
  }
  if (section === SectionEmail) {
    return <Notification />;
  }
  if (section === SectionBilling) {
    return <Billing />;
  }

  return <div>Not found</div>;
}

function Settings({ match }) {
  const { params } = match;
  const { section } = params;

  return (
    <div className="page">
      <Helmet>
        <meta name="description" content="Dnote settings" />
      </Helmet>

      {renderContent(section)}
    </div>
  );
}

export default Settings;
