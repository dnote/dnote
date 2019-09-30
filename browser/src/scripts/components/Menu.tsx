import React, { Fragment } from 'react';

import Link from './Link';

export default ({ toggleMenu, loggedIn, onLogout }) => (
  <Fragment>
    <ul className="menu">
      <li>
        <Link to="/" onClick={toggleMenu} className="menu-link">
          Home
        </Link>
      </li>

      {loggedIn && (
        <li>
          <form
            onSubmit={e => {
              e.preventDefault();

              onLogout(toggleMenu);
            }}
          >
            <input
              type="submit"
              value="Logout"
              className="menu-link logout-button"
            />
          </form>
        </li>
      )}
    </ul>

    <div className="menu-overlay" onClick={toggleMenu} />
  </Fragment>
);
