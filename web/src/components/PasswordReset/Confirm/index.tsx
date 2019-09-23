import React, { useState } from 'react';
import Helmet from 'react-helmet';
import { Link, withRouter, RouteComponentProps } from 'react-router-dom';

import { homePathDef } from 'web/libs/paths';
import services from 'web/libs/services';
import Form from './Form';
import Logo from '../../Icons/Logo';
import { getCurrentUser } from '../../../store/auth';
import { setMessage } from '../../../store/ui';
import { useDispatch } from '../../../store';
import authStyles from '../../Common/Auth.scss';
import Flash from '../../Common/Flash';

interface Match {
  token: string;
}

interface Props extends RouteComponentProps<Match> {}

const PasswordResetConfirm: React.SFC<Props> = ({ match, history }) => {
  const [errorMsg, setErrorMsg] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const dispatch = useDispatch();

  function handleResetPassword(password: string, passwordConfirmation: string) {
    if (!password) {
      setErrorMsg('Please enter password');
      return;
    }
    if (!passwordConfirmation) {
      setErrorMsg('Please enter password confirmation');
      return;
    }

    const { token } = match.params;

    setSubmitting(true);
    setErrorMsg('');

    services.users
      .resetPassword({ token, password })
      .then(() => {
        return dispatch(getCurrentUser());
      })
      .then(() => {
        dispatch(
          setMessage({
            message: 'Your password was successfully reset.',
            kind: 'info',
            path: homePathDef
          })
        );
        history.push('/');
      })
      .catch(err => {
        setSubmitting(false);
        setErrorMsg(err.message);
      });
  }

  return (
    <div className={authStyles.page}>
      <Helmet>
        <title>Reset Password</title>
      </Helmet>
      <div className="container">
        <Link to="/">
          <Logo fill="#252833" width={60} height={60} />
        </Link>
        <h1 className={authStyles.heading}>Reset Password</h1>

        <div className={authStyles.body}>
          <Flash kind="info">
            Password must be at least 8 characters long.
          </Flash>

          <div className={authStyles.panel}>
            {errorMsg && <div className="alert alert-danger">{errorMsg}</div>}

            <Form
              onResetPassword={handleResetPassword}
              submitting={submitting}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default withRouter(PasswordResetConfirm);
