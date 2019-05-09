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

import { getLessonCalendar } from '../../actions/calendar';

import Header from '../Common/Page/Header';
import Body from '../Common/Page/Body';
import SubscriberWall from '../Common/SubscriberWall';
import Content from './Content';

function Dashboard({ demo, userData, calendarData, doGetLessonCalendar }) {
  useEffect(() => {
    doGetLessonCalendar({ demo });
  }, [demo, doGetLessonCalendar]);

  const user = userData.data;

  return (
    <div className="dashboard-page page">
      <Helmet>
        <title>Dashboard</title>
      </Helmet>

      <Header heading="Dashboard" />

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
  doGetLessonCalendar: getLessonCalendar
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Dashboard);
