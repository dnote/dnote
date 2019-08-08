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
import { connect } from 'react-redux';

import { receiveUser } from '../../store/auth';
import { updateMessage } from '../../store/ui';
import * as usersService from '../../services/users';

function VerifyEmail({ match, history, doReceiveUser, doUpdateMessage }) {
  const { token } = match.params;

  useEffect(() => {
    usersService
      .verifyEmail({ token })
      .then(res => {
        doReceiveUser(res);
        doUpdateMessage('Email was successfully verified', 'info');
        history.push('/');
      })
      .catch(err => {
        doUpdateMessage(err.message, 'error');
        history.push('/');
      });
  }, [doReceiveUser, doUpdateMessage, history, token]);

  return <div />;
}

const mapDispatchToProps = {
  doReceiveUser: receiveUser,
  doUpdateMessage: updateMessage
};

export default connect(
  null,
  mapDispatchToProps
)(VerifyEmail);
