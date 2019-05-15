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
import { joinPath, isNotePath, isDemoPath } from './libs/paths';

// Components
import Home from './components/Home';
import Books from './components/Books';
import Login from './components/Login';
import Digests from './components/Digests';
import Join from './components/Join';
import Settings from './components/Settings';
import NotFound from './components/Common/NotFound';
import VerifyEmail from './components/VerifyEmail';
import EmailPreference from './components/EmailPreference';
import Note from './components/Note';
import Digest from './components/Digest';
import Subscription from './components/Subscription';

import LegacyLogin from './components/LegacyLogin';
import LegacyJoin from './components/LegacyJoin';
import LegacyEncrypt from './components/LegacyEncrypt';

const AuthenticatedHome = userOnly(Home);
const AuthenticatedBooks = userOnly(Books);
const AuthenticatedSettings = userOnly(Settings);
const AuthenticatedDigest = userOnly(Digest);
const AuthenticatedDigests = userOnly(Digests);
const AuthenticatedNote = userOnly(Note);
const AuthenticatedSubscription = userOnly(Subscription, joinPath().pathname);
const GuestLogin = guestOnly(Login);
const GuestJoin = guestOnly(Join);

export default function render(isEditor) {
  const routes = [
    {
      path: ['/', '/notes/:noteUUID', '/demo/notes/:noteUUID'],
      exact: true,
      render: ({ location }) => {
        const demo = isDemoPath(location.pathname);

        if (isNotePath(location.pathname) && !isEditor) {
          if (demo) {
            return <Note demo />;
          }

          return <AuthenticatedNote />;
        }

        if (demo) {
          return <Home demo />;
        }

        return <AuthenticatedHome />;
      }
    },
    {
      path: '/demo',
      exact: true,
      render: () => {
        return <Home demo />;
      }
    },
    {
      path: '/demo/books',
      exact: true,
      render: () => {
        return <Books demo />;
      }
    },
    {
      path: '/demo/digests',
      exact: true,
      render: () => {
        return <Digests demo />;
      }
    },
    {
      path: '/demo/notes/:noteUUID',
      exact: true,
      render: () => {
        return <Note demo />;
      }
    },
    {
      path: '/login',
      exact: true,
      component: GuestLogin
    },
    {
      path: '/join',
      exact: true,
      component: GuestJoin
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
      path: '/books',
      exact: true,
      component: AuthenticatedBooks
    },
    {
      path: '/digests',
      exact: true,
      component: AuthenticatedDigests
    },
    {
      path: '/settings/:section',
      exact: true,
      component: AuthenticatedSettings
    },
    {
      path: '/digests/:digestUUID',
      exact: true,
      component: AuthenticatedDigest
    },
    {
      path: '/demo/digests/:digestUUID',
      exact: true,
      render: () => {
        return <Digest demo />;
      }
    },
    {
      path: '/subscriptions',
      exact: true,
      component: AuthenticatedSubscription
    },
    {
      path: '/legacy/login',
      exact: true,
      component: LegacyLogin
    },
    {
      path: '/legacy/register',
      exact: true,
      component: LegacyJoin
    },
    {
      path: '/legacy/encrypt',
      exact: true,
      component: LegacyEncrypt
    },
    {
      component: NotFound
    }
  ];

  return renderRoutes(routes);
}
