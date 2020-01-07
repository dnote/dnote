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

import React from 'react';

import Link from './Link';
import MenuToggleIcon from './MenuToggleIcon';
import CloseIcon from './CloseIcon';

interface Props {
  toggleMenu: () => void;
  isShowingMenu: boolean;
}

const Header: React.FunctionComponent<Props> = ({
  toggleMenu,
  isShowingMenu
}) => (
  <header className="header">
    <Link to="/" className="logo-link" tabIndex={-1}>
      <img src="images/logo-circle.png" alt="dnote" className="logo" />
    </Link>

    <a
      href="#toggle"
      className="menu-toggle"
      onClick={e => {
        e.preventDefault();

        toggleMenu();
      }}
      tabIndex={-1}
    >
      {isShowingMenu ? <CloseIcon /> : <MenuToggleIcon />}
    </a>
  </header>
);

export default Header;
