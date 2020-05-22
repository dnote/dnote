/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
import Settings from './components/Settings';
import NotFound from './components/Common/NotFound';
import VerifyEmail from './components/VerifyEmail';
import New from './components/New';
import Edit from './components/Edit';
import Note from './components/Note';
import Books from './components/Books';
import PasswordResetRequest from './components/PasswordReset/Request';
import PasswordResetConfirm from './components/PasswordReset/Confirm';
import EmailPreference from './components/EmailPreference';

// paths
import {
  notePathDef,
  homePathDef,
  booksPathDef,
  loginPathDef,
  joinPathDef,
  noteEditPathDef,
  noteNewPathDef,
  settingsPathDef,
  passwordResetRequestPathDef,
  passwordResetConfirmPathDef,
  verifyEmailPathDef,
  emailPreferencePathDef
} from './libs/paths';

const AuthenticatedHome = userOnly(Home);
const AuthenticatedNew = userOnly(New);
const AuthenticatedEdit = userOnly(Edit);
const AuthenticatedBooks = userOnly(Books);
const GuestJoin = guestOnly(Join);
const GuestLogin = guestOnly(Login);
const GuestPasswordResetRequest = guestOnly(PasswordResetRequest);
const GuestPasswordResetConfirm = guestOnly(PasswordResetConfirm);
const AuthenticatedSettings = userOnly(Settings);

const routes = [
  {
    path: homePathDef,
    exact: true,
    component: AuthenticatedHome
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
    path: settingsPathDef,
    exact: true,
    component: AuthenticatedSettings
  },
  {
    path: verifyEmailPathDef,
    exact: true,
    component: VerifyEmail
  },
  {
    path: noteNewPathDef,
    exact: true,
    component: AuthenticatedNew
  },
  {
    path: passwordResetRequestPathDef,
    exact: true,
    component: GuestPasswordResetRequest
  },
  {
    path: passwordResetConfirmPathDef,
    exact: true,
    component: GuestPasswordResetConfirm
  },
  {
    path: emailPreferencePathDef,
    exact: true,
    component: EmailPreference
  },
  {
    component: NotFound
  }
];

export default function render(): React.ReactNode {
  return renderRoutes(routes);
}
