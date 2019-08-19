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
import { renderRoutes } from 'react-router-config';

import userOnly from './hocs/userOnly';
import guestOnly from './hocs/guestOnly';

// Components
import Home from './components/Home';
import Login from './components/Login';
import Join from './components/Join';
import NotFound from './components/Common/NotFound';
import VerifyEmail from './components/VerifyEmail';
import EmailPreference from './components/EmailPreference';
import New from './components/New';
import Edit from './components/Edit';
import Note from './components/Note';
import Books from './components/Books';
import {
  notePathDef,
  homePathDef,
  booksPathDef,
  loginPathDef,
  joinPathDef,
  noteEditPathDef,
  noteNewPathDef
} from './libs/paths';

const AuthenticatedHome = userOnly(Home);
const AuthenticatedNew = userOnly(New);
const AuthenticatedEdit = userOnly(Edit);
const AuthenticatedBooks = userOnly(Books);
const GuestJoin = guestOnly(Join);
const GuestLogin = guestOnly(Login);

const routes = [
  {
    path: homePathDef,
    exact: true,
    render: () => {
      return <AuthenticatedHome />;
    }
  },
  {
    path: loginPathDef,
    exact: true,
    component: GuestLogin
  },
  {
    path: joinPathDef,
    exact: true,
    component: GuestJoin
  },
  {
    path: notePathDef,
    exact: true,
    component: Note
  },
  {
    path: booksPathDef,
    exact: true,
    component: AuthenticatedBooks
  },
  {
    path: noteEditPathDef,
    exact: true,
    component: AuthenticatedEdit
  },
  {
    path: '/verify-email/:token',
    exact: true,
    component: VerifyEmail
  },
  {
    path: '/email-preference',
    exact: true,
    component: EmailPreference
  },
  {
    path: noteNewPathDef,
    exact: true,
    component: AuthenticatedNew
  },
  {
    component: NotFound
  }
];

export default function render(): React.ReactNode {
  return renderRoutes(routes);
}
