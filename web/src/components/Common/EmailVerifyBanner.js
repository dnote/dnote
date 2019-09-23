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

import React, { Fragment, useState } from 'react';
import { connect } from 'react-redux';
import classnames from 'classnames';

import Flash from './Flash';
import Button from './Button';

import * as usersService from 'jslib/services/users';

import styles from './EmailVerifyBanner.module.scss';

function EmailVerifyBanner({ demo, userData }) {
  const [submitted, setSubmitted] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  const user = userData.data;

  if (demo || !user.cloud || user.email_verified) {
    return null;
  }

  function handleSubmit(e) {
    e.preventDefault();

    setSubmitting(true);
    setErrorMsg('');

    usersService
      .sendEmailVerificationEmail(user.api_key)
      .then(() => {
        setSubmitting(false);
        setSubmitted(true);
      })
      .catch(err => {
        setSubmitting(false);
        setErrorMsg(err.message);
      });
  }

  function renderMessage() {
    const { email } = user;

    if (submitted) {
      return (
        <Fragment>
          Your verification email was sent to
          <strong className={styles.address}>{email}.</strong> If you do not
          receive it, please check spam.
        </Fragment>
      );
    }

    if (errorMsg) {
      return (
        <Fragment>
          Error sending verification email to{' '}
          <strong className={styles.address}>{email}</strong>: {errorMsg}
        </Fragment>
      );
    }

    return (
      <Fragment>
        You need to verify your email
        <strong className={styles.address}>{email}</strong>
        to receive weekly digest.
      </Fragment>
    );
  }

  let type;

  if (errorMsg) {
    type = 'danger';
  } else {
    type = 'info';
  }

  return (
    <Flash
      id="T-email-verify-banner"
      wrapperClassName={classnames(styles.banner, {
        'T-submitted': submitted
      })}
      contentClassName={styles.wrapper}
      type={type}
    >
      <div>{renderMessage()}</div>

      {!submitted && (
        <form onSubmit={handleSubmit} className={styles.form}>
          <Button type="submit" kind="second" isBusy={submitting}>
            Send verification email
          </Button>
        </form>
      )}
    </Flash>
  );
}

function mapStateToProps(state) {
  return {
    userData: state.auth.user
  };
}

export default connect(mapStateToProps)(EmailVerifyBanner);
