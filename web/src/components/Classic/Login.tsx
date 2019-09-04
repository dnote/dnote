import React, { useState } from 'react';
import Helmet from 'react-helmet';

import { useDispatch } from '../../store';
import { loginHelper, aes256GcmDecrypt } from '../../crypto';
import { bufToB64, b64ToBuf } from '../../libs/encoding';
import { getCurrentUser } from '../../store/auth';
import { classicPresignin, classicSignin } from '../../services/users';
import authStyles from '../Common/Auth.scss';
import Logo from '../Icons/Logo';
import Flash from '../Common/Flash';
import LoginForm from '../Login/LoginForm';

interface Props {}

const ClassicLogin: React.SFC<Props> = () => {
  const [errMsg, setErrMsg] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [emailVal, setEmailVal] = useState('');
  const dispatch = useDispatch();

  async function handleLogin(email, password) {
    if (!email) {
      setErrMsg('Please enter email');
      return;
    }
    if (!password) {
      setErrMsg('Please enter password');
      return;
    }

    setErrMsg('');
    setSubmitting(true);

    try {
      const { iteration } = await classicPresignin({ email });

      if (iteration === 0) {
        throw new Error('Please login from /login');
      }

      const { masterKey, authKey } = await loginHelper({
        email,
        password,
        iteration
      });
      const signinResp = await classicSignin({ email, authKey });

      const cipherKey = await aes256GcmDecrypt(
        b64ToBuf(masterKey),
        b64ToBuf(signinResp.cipherKeyEnc)
      );
      localStorage.setItem('cipherKey', bufToB64(cipherKey));

      await dispatch(getCurrentUser());
    } catch (err) {
      console.log(err);
      setErrMsg(err.message);
      setSubmitting(false);
    }
  }

  return (
    <div className={authStyles.page}>
      <Helmet>
        <title>Sign In (Classic)</title>
      </Helmet>
      <div className="container">
        <a href="/">
          <Logo fill="#252833" width={60} height={60} />
        </a>
        <h1 className={authStyles.heading}>Sign in to Dnote Classic</h1>
        <div className={authStyles.body}>
          <div className={authStyles.panel}>
            {errMsg && (
              <Flash
                id="T-login-error"
                kind="danger"
                wrapperClassName={authStyles['error-flash']}
              >
                {errMsg}
              </Flash>
            )}

            <LoginForm
              onLogin={handleLogin}
              submitting={submitting}
              email={emailVal}
              onUpdateEmail={val => setEmailVal(val)}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default ClassicLogin;
