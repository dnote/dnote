/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import React, { useState, useEffect } from 'react';
import classnames from 'classnames';

import initServices from '../utils/services';
import { resetSettings } from '../store/settings/actions';
import { SettingsState } from '../store/settings/types';
import { useSelector, useDispatch } from '../store/hooks';
import Header from './Header';
import Home from './Home';
import Menu from './Menu';
import Success from './Success';
import Settings from './Settings';
import Composer from './Composer';

interface Props {}

function renderRoutes(path: string, isLoggedIn: boolean) {
  switch (path) {
    case '/success':
      return <Success />;
    case '/': {
      if (isLoggedIn) {
        return <Composer />;
      }

      return <Home />;
    }
    case '/settings': {
      return <Settings />;
    }
    default:
      return <div>Not found</div>;
  }
}

// useCheckSessionValid ensures that the current session is valid
function useCheckSessionValid(settings: SettingsState) {
  const dispatch = useDispatch();

  useEffect(() => {
    // if session is expired, clear it
    const now = Math.round(new Date().getTime() / 1000);
    if (settings.sessionKey && settings.sessionKeyExpiry < now) {
      dispatch(resetSettings());
    }
  }, [dispatch, settings.sessionKey, settings.sessionKeyExpiry]);
}

const App: React.FunctionComponent<Props> = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [errMsg, setErrMsg] = useState('');

  const dispatch = useDispatch();
  const { path, settings } = useSelector(state => {
    return {
      path: state.location.path,
      settings: state.settings
    };
  });

  useCheckSessionValid(settings);

  const isLoggedIn = Boolean(settings.sessionKey);
  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  const handleLogout = async (done?: Function) => {
    try {
      await initServices(settings.apiUrl).users.signout();
      dispatch(resetSettings());

      if (done) {
        done();
      }
    } catch (e) {
      setErrMsg(e.message);
    }
  };

  return (
    <div className="container">
      <Header toggleMenu={toggleMenu} isShowingMenu={isMenuOpen} />

      {isMenuOpen && (
        <Menu
          toggleMenu={toggleMenu}
          loggedIn={isLoggedIn}
          onLogout={handleLogout}
        />
      )}

      <main>
        {errMsg && <div className="alert error">{errMsg}</div>}

        {renderRoutes(path, isLoggedIn)}
      </main>
    </div>
  );
};

export default App;
