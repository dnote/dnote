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

import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import Helmet from 'react-helmet';

import { getDigests } from '../../actions/digests';

import Header from '../Common/Page/Header';
import Body from '../Common/Page/Body';
import SubscriberWall from '../Common/SubscriberWall';
import Content from './Content';

function Digests({ demo, userData, calendarData, doGetDigests }) {
  useEffect(() => {
    doGetDigests({ demo });
  }, [demo, doGetDigests]);

  const user = userData.data;

  return (
    <div className="dashboard-page page">
      <Helmet>
        <title>Digests</title>
      </Helmet>

      <Header heading="Digests" />

      <Body>
        <div className="container">
          {demo || user.cloud ? (
            <Content user={user} calendar={calendarData} demo={demo} />
          ) : (
            <SubscriberWall />
          )}
        </div>
      </Body>
    </div>
  );
}

function mapStateToProps(state) {
  return {
    userData: state.auth.user,
    calendarData: state.calendar
  };
}

const mapDispatchToProps = {
  doGetDigests: getDigests
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Digests);
