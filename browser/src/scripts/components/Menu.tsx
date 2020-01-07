/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
      <li>
        <Link to="/settings" onClick={toggleMenu} className="menu-link">
          Settings
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

    <div
      className="menu-overlay"
      onClick={toggleMenu}
      onKeyDown={() => {}}
      role="none"
    />
  </Fragment>
);
