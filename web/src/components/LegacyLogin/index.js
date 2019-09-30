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
import Helmet from 'react-helmet';
import { Link, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import { getReferrer } from 'jslib/helpers/url';
import OauthLoginButton from './OauthLoginButton';
import LoginForm from './LoginForm';
import LegacyFooter from '../Common/LegacyFooter';

import Logo from '../Icons/Logo';
import { legacySignin } from 'jslib/services/users';
import { receiveUser } from '../../actions/auth';
import { updateAuthEmail } from '../../actions/form';

import './module.scss';

class LegacyLogin extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMsg: '',
      submitting: false
    };
  }

  handlePasswordLogin = (email, password) => {
    if (!email) {
      this.setState({ errorMsg: 'Please enter email' });
      return;
    }
    if (!password) {
      this.setState({ errorMsg: 'Please enter password' });
      return;
    }

    this.setState({ submitting: true, errorMsg: '' }, () => {
      legacySignin({ email, password })
        .then(res => {
          const { history, doReceiveUser } = this.props;
          const { user } = res;

          doReceiveUser(user);

          history.push('/legacy/register');
        })
        .catch(err => {
          this.setState({ submitting: false, errorMsg: err.message });
        });
    });
  };

  render() {
    const { location, doUpdateAuthFormEmail, email } = this.props;
    const { submitting, errorMsg } = this.state;

    const referrer = getReferrer(location);

    return (
      <div className="auth-page login-page">
        <Helmet>
          <title>Legacy Login</title>
        </Helmet>
        <div className="container">
          <Link to="/">
            <Logo fill="#252833" width="60" height="60" />
          </Link>
          <h1 className="heading">Sign into new Dnote</h1>

          <div className="auth-body">
            <div className="auth-panel">
              <OauthLoginButton
                referrer={referrer}
                provider="github"
                text="Sign in with GitHub"
              />
              <OauthLoginButton
                referrer={referrer}
                provider="gplus"
                text="Sign in with Google"
              />

              <div className="divider-text">or</div>

              {errorMsg && <div className="alert alert-danger">{errorMsg}</div>}

              <LoginForm
                email={email}
                onLogin={this.handlePasswordLogin}
                submitting={submitting}
                onEmailChange={doUpdateAuthFormEmail}
              />
            </div>
          </div>

          <LegacyFooter />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    email: state.form.auth.email
  };
}

const mapDispatchToProps = {
  doUpdateAuthFormEmail: updateAuthEmail,
  doReceiveUser: receiveUser
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(LegacyLogin)
);
