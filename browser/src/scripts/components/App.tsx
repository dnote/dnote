import React, { useState, useEffect } from 'react';
import classnames from 'classnames';

import services from '../utils/services';
import { resetSettings } from '../store/settings/actions';
import { useSelector, useDispatch } from '../store/hooks';
import Header from './Header';
import Home from './Home';
import Menu from './Menu';
import Success from './Success';
import Composer from './Composer';

interface Props {}

function renderRoutes(path: string, isLoggedIn: boolean) {
  switch (path) {
    case '/success':
      return <Success />;
    case '/':
      if (isLoggedIn) {
        return <Composer />;
      }

      return <Home />;
    default:
      return <div>Not found</div>;
  }
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

  useEffect(() => {
    // if session is expired, clear it
    const now = Math.round(new Date().getTime() / 1000);
    if (settings.sessionKey && settings.sessionKeyExpiry < now) {
      dispatch(resetSettings());
    }
  }, [dispatch]);

  const isLoggedIn = Boolean(settings.sessionKey);
  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };
  const handleLogout = async (done?: Function) => {
    try {
      await services.users.signout();
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
