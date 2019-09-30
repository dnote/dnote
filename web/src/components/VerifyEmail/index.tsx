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

import React, { useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import { getHomePath, homePathDef, emailPrefPathDef } from 'web/libs/paths';
import services from 'web/libs/services';
import { receiveUser } from '../../store/auth';
import { setMessage } from '../../store/ui';
import { useDispatch } from '../../store';

interface Match {
  token: string;
}

interface Props extends RouteComponentProps<Match> {}

const VerifyEmail: React.SFC<Props> = ({ match, history }) => {
  const dispatch = useDispatch();
  const { token } = match.params;

  useEffect(() => {
    const homePath = getHomePath();

    services.users
      .verifyEmail({ token })
      .then(res => {
        dispatch(receiveUser(res));
        dispatch(
          setMessage({
            message: 'Email was successfully verified',
            kind: 'info',
            path: homePathDef
          })
        );
        history.push(homePath);
      })
      .catch(err => {
        dispatch(
          setMessage({
            message: err.message,
            kind: 'error',
            path: emailPrefPathDef
          })
        );
      });
  }, [dispatch, history, token]);

  return null;
};

export default VerifyEmail;
