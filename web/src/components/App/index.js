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

import { hot } from 'react-hot-loader/root';
import React, { useEffect, useRef } from 'react';
import { connect } from 'react-redux';
import { Switch, Route } from 'react-router';
import classnames from 'classnames';
import isEqual from 'lodash/isEqual';
import { withRouter } from 'react-router-dom';

import DemoHeader from '../Header/DemoHeader';
import NoteHeader from '../Header/NoteHeader';
import SubscriptionsHeader from '../Header/SubscriptionsHeader';
import HeaderData from './HeaderData';
import NoteSidebar from '../Home/NoteSidebar';
import MainSidebar from '../Common/Sidebar/MainSidebar';
import Flash from '../Common/Flash';
import SettingsSidebar from '../Common/Sidebar/SettingsSidebar';
import Splash from '../Splash';
import Footer from '../Footer';
import MainFooter from '../Footer/MainFooter';
import { getCurrentUser } from '../../actions/auth';
import render from '../../routes';
import { footerPaths, isDemoPath, checkBoxedLayout } from '../../libs/paths';
import SystemMessage from '../Common/SystemMessage';

import './App.global.scss';
import styles from './App.module.scss';

function checkIsEditor(location, prevLocation) {
  return (
    location.state &&
    location.state.editor &&
    // not initial render
    !isEqual(prevLocation, location)
  );
}

function App({ location, history, doGetCurrentUser, user }) {
  const prevLocationRef = useRef(location);

  useEffect(() => {
    doGetCurrentUser();
  }, [doGetCurrentUser]);

  useEffect(() => {
    if (
      history.action !== 'POP' &&
      (!location.state || !location.state.editor)
    ) {
      prevLocationRef.current = location;
    }
  });

  if (user.errorMessage) {
    return (
      <Flash type="danger">
        Could not fetch the current user: <span>{user.errorMessage}</span>
      </Flash>
    );
  }

  const isReady = !user.isFetched || user.isFetching;

  if (isReady) {
    return <Splash />;
  }

  const isDemo = isDemoPath(location.pathname);
  const isEditor = checkIsEditor(location, prevLocationRef.current);
  const isBoxedLayout = checkBoxedLayout(location, isEditor);

  return (
    <div
      className={classnames(styles.app, {
        [styles.boxed]: isBoxedLayout,
        [styles.full]: !isBoxedLayout
      })}
    >
      <HeaderData />
      <SystemMessage />

      <Route path="/demo" component={DemoHeader} />
      <Route
        path={[
          '/notes/:noteUUID',
          '/demo/notes/:noteUUID',
          '/digests/:digestUUID',
          '/demo/digests/:digestUUID'
        ]}
        render={() => {
          if (isEditor) {
            return null;
          }

          return <NoteHeader demo={isDemo} />;
        }}
      />
      <Route path="/subscriptions" component={SubscriptionsHeader} />

      <main
        className={classnames(styles.wrapper, {
          [styles['boxed-wrapper']]: isBoxedLayout,
          [styles['full-wrapper']]: !isBoxedLayout
        })}
      >
        <MainSidebar isEditor={isEditor} demo={isDemo} />
        <NoteSidebar isEditor={isEditor} demo={isDemo} />
        <Route
          path="/settings/:section"
          render={() => {
            return <SettingsSidebar />;
          }}
        />

        <div
          className={classnames(styles.main, {
            [styles['full-main']]: !isBoxedLayout
          })}
        >
          <Switch>{render(isEditor)}</Switch>
        </div>
      </main>
      <Route
        exact
        path={footerPaths}
        render={() => {
          return <Footer isEditor={isEditor} demo={isDemo} />;
        }}
      />
      <Route
        path="/subscriptions"
        render={() => {
          return <MainFooter />;
        }}
      />
    </div>
  );
}

function mapStateToProps(state) {
  return {
    user: state.auth.user,
    message: state.ui.message,
    layout: state.ui.layout
  };
}

const mapDispatchToProps = {
  doGetCurrentUser: getCurrentUser
};

export default hot(
  withRouter(
    connect(
      mapStateToProps,
      mapDispatchToProps
    )(App)
  )
);
