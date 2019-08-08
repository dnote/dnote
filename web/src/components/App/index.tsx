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

import React, { useEffect, Fragment } from 'react';
import { hot } from 'react-hot-loader/root';
import { Switch, Route } from 'react-router';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import { Location } from 'history';

import Splash from '../Splash';
import { getCurrentUser } from '../../store/auth';
import { getBooks } from '../../store/books';
import { setPrevLocation } from '../../store/route';
import { useDispatch, useSelector } from '../../store';
import { usePrevious } from '../../libs/hooks';
import HeaderData from './HeaderData';
import render from '../../routes';
import NoteHeader from '../Header/Note';
import NormalHeader from '../Header/Normal';
import TabBar from '../TabBar';
import SystemMessage from '../Common/SystemMessage';
import styles from './App.scss';

import './App.global.scss';

interface Props extends RouteComponentProps<any> {}

function useFetchData() {
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(getCurrentUser()).then(u => {
      if (!u) {
        return null;
      }

      return dispatch(getBooks());
    });
  }, [dispatch]);
}

function hasLocationChanged(loc1: Location<any>, loc2: Location<any>) {
  return (
    loc1.pathname !== loc2.pathname ||
    loc1.search !== loc2.search ||
    loc1.hash !== loc2.hash
  );
}

// useSavePrevLocation saves the prev location upon route change
function useSavePrevLocation(location: Location) {
  const prevLocation = usePrevious(location);
  const dispatch = useDispatch();

  useEffect(() => {
    if (!prevLocation) {
      return;
    }

    if (hasLocationChanged(location, prevLocation)) {
      console.log(location, prevLocation);
      dispatch(setPrevLocation(prevLocation));
    }
  }, [prevLocation, dispatch, location]);
}

const App: React.SFC<Props> = ({ location }) => {
  useFetchData();
  useSavePrevLocation(location);

  const { user } = useSelector(state => {
    return {
      user: state.auth.user
    };
  });

  const isReady = !user.isFetched || user.isFetching;
  if (isReady) {
    return <Splash />;
  }

  return (
    <Fragment>
      <HeaderData />

      <Switch>
        <Route path={['/login', '/join']} exact component={null} />
        <Route path="/notes/:noteUUID" exact render={() => <NoteHeader />} />
        <Route path="/" component={NormalHeader} />
      </Switch>

      <main className={styles.wrapper}>
        <SystemMessage />
        <Switch>{render()}</Switch>
      </main>

      <Switch>
        <Route path={['/login', '/join']} exact component={null} />
        <Route path="/" component={TabBar} />
      </Switch>
    </Fragment>
  );
};

export default hot(withRouter(App));
