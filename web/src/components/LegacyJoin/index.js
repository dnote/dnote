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
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import JoinForm from './JoinForm';
import Logo from '../Icons/Logo';

import { legacyRegister } from 'jslib/services/users';
import { receiveUser, legacyGetCurrentUser } from '../../actions/auth';
import { updateAuthEmail } from '../../actions/form';
import LegacyFooter from '../Common/LegacyFooter';
import { registerHelper } from '../../crypto';
import { DEFAULT_KDF_ITERATION } from '../../crypto/consts';

class LegacyJoin extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMsg: '',
      submitting: false,
      isReady: false
    };
  }

  componentDidMount() {
    const { doLegacyGetCurrentUser } = this.props;

    doLegacyGetCurrentUser().then(() => {
      const { userData, history } = this.props;

      const user = userData.data;

      if (!user.legacy) {
        history.push('/legacy/encrypt');
      } else {
        this.setState({ isReady: true });
      }
    });
  }

  handleJoin = (email, password, passwordConfirmation) => {
    const { history } = this.props;

    if (!email) {
      this.setState({ errorMsg: 'Please enter email' });
      return;
    }
    if (!password) {
      this.setState({ errorMsg: 'Please enter password' });
      return;
    }
    if (!passwordConfirmation) {
      this.setState({ errorMsg: 'The passwords do not match' });
      return;
    }

    this.setState({ submitting: true, errorMsg: '' }, async () => {
      try {
        const { cipherKey, cipherKeyEnc, authKey } = await registerHelper({
          email,
          password,
          iteration: DEFAULT_KDF_ITERATION
        });
        await legacyRegister({
          email,
          authKey,
          cipherKeyEnc,
          iteration: DEFAULT_KDF_ITERATION
        });

        localStorage.setItem('cipherKey', cipherKey);

        history.push('/legacy/encrypt');
      } catch (err) {
        console.log(err);
        this.setState({ submitting: false, errorMsg: err.message });
      }
    });
  };

  render() {
    const { doUpdateAuthFormEmail, email } = this.props;
    const { errorMsg, submitting, isReady } = this.state;

    if (!isReady) {
      return <div>Loading...</div>;
    }

    return (
      <div className="auth-page">
        <Helmet>
          <title>Join</title>
        </Helmet>

        <div className="container">
          <div className="container">
            <a href="/">
              <Logo fill="#252833" width="60" height="60" />
            </a>
            <h1 className="heading">Choose your email and password</h1>

            <div className="auth-body">
              <div className="auth-panel">
                {errorMsg && (
                  <div className="alert alert-danger">{errorMsg}</div>
                )}

                <JoinForm
                  onJoin={this.handleJoin}
                  submitting={submitting}
                  onEmailChange={doUpdateAuthFormEmail}
                  email={email}
                />
              </div>
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
    email: state.form.auth.email,
    userData: state.auth.user
  };
}

const mapDispatchToProps = {
  doUpdateAuthFormEmail: updateAuthEmail,
  doReceiveUser: receiveUser,
  doLegacyGetCurrentUser: legacyGetCurrentUser
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(LegacyJoin)
);
