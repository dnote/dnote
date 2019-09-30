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

import React, { useState } from 'react';
import Helmet from 'react-helmet';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import services from 'web/libs/services';
import { getClassicMigrationPath, ClassicMigrationSteps } from 'web/libs/paths';
import Logo from '../Icons/Logo';
import authStyles from '../Common/Auth.scss';
import Flash from '../Common/Flash';
import JoinForm from '../Join/JoinForm';

interface Props extends RouteComponentProps {}

const SetPassword: React.SFC<Props> = ({ history }) => {
  const [errMsg, setErrMsg] = useState('');
  const [submitting, setSubmitting] = useState(false);

  async function handleJoin(email, password, passwordConfirmation) {
    if (!password) {
      setErrMsg('Please enter a pasasword.');
      return;
    }
    if (!passwordConfirmation) {
      setErrMsg('The passwords do not match.');
      return;
    }

    setErrMsg('');
    setSubmitting(true);

    try {
      await services.users.classicSetPassword({ password });

      const p = getClassicMigrationPath(ClassicMigrationSteps.decrypt);
      history.push(p);
    } catch (err) {
      console.log(err);
      setErrMsg(err.message);
      setSubmitting(false);
    }
  }

  return (
    <div className={authStyles.page}>
      <Helmet>
        <title>Set password (Classic)</title>
      </Helmet>

      <div className="container">
        <a href="/">
          <Logo fill="#252833" width={60} height={60} />
        </a>
        <h1 className={authStyles.heading}>Set password</h1>

        <div className={authStyles.body}>
          <div className={authStyles.panel}>
            {errMsg && (
              <Flash kind="danger" wrapperClassName={authStyles['error-flash']}>
                {errMsg}
              </Flash>
            )}

            <JoinForm
              onJoin={handleJoin}
              submitting={submitting}
              cta="Confirm"
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default withRouter(SetPassword);
