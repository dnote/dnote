import React, { useState } from 'react';
import Helmet from 'react-helmet';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import Logo from '../Icons/Logo';
import authStyles from '../Common/Auth.scss';
import Flash from '../Common/Flash';
import JoinForm from '../Join/JoinForm';
import * as usersService from '../../services/users';
import {
  getClassicMigrationPath,
  ClassicMigrationSteps
} from '../../libs/paths';

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
      await usersService.classicSetPassword({ password });

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
