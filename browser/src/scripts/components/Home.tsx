import React, { useState } from 'react';
import { connect } from 'react-redux';
import { findDOMNode } from 'react-dom';

import Link from './Link';
import config from '../utils/config';
import { updateSettings } from '../store/settings/actions';
import { useDispatch } from '../store/hooks';
import services from '../utils/services';

interface Props {}

const Home: React.FunctionComponent<Props> = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errMsg, setErrMsg] = useState('');
  const [loggingIn, setLoggingIn] = useState(false);
  const dispatch = useDispatch();

  const handleLogin = async e => {
    e.preventDefault();

    setErrMsg('');
    setLoggingIn(true);

    try {
      const signinResp = await services.users.signin({ email, password });

      dispatch(
        updateSettings({
          sessionKey: signinResp.key,
          sessionKeyExpiry: signinResp.expiresAt
        })
      );
    } catch (e) {
      console.log('error while logging in', e);

      setErrMsg(e.message);
      setLoggingIn(false);
    }
  };

  return (
    <div className="home">
      <h1 className="greet">Welcome to Dnote</h1>

      <p className="lead">A simple personal knowledge base</p>

      {errMsg && <div className="alert error">{errMsg}</div>}

      <form id="login-form" onSubmit={handleLogin}>
        <label htmlFor="email-input">Email</label>

        <input
          type="email"
          placeholder="your@email.com"
          className="input login-input"
          id="email-input"
          value={email}
          onChange={e => {
            setEmail(e.target.value);
          }}
        />

        <label htmlFor="password-input">Password</label>
        <input
          type="password"
          placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;"
          className="input login-input"
          id="password-input"
          value={password}
          onChange={e => {
            setPassword(e.target.value);
          }}
        />

        <button
          type="submit"
          className="button button-first button-small login-btn"
          disabled={loggingIn}
        >
          {loggingIn ? 'Signing in...' : 'Signin'}
        </button>
      </form>

      <div className="actions">
        Don't have an account?{' '}
        <a
          href="https://app.getdnote.com/join"
          target="_blank"
          className="signup"
        >
          Sign Up
        </a>
      </div>
    </div>
  );
};

export default Home;
