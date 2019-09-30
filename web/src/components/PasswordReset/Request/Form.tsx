import React, { useState } from 'react';
import classnames from 'classnames';

import Button from '../../Common/Button';
import authStyles from '../../Common/Auth.scss';
import styles from './Form.scss';

interface Props {
  onSubmit: (email: string) => void;
  submitting: boolean;
}

const PasswordResetRequestForm: React.SFC<Props> = ({
  onSubmit,
  submitting
}) => {
  const [email, setEmail] = useState('');

  return (
    <form
      onSubmit={e => {
        e.preventDefault();

        onSubmit(email);
      }}
      className={authStyles.form}
    >
      <div className={authStyles['input-row']}>
        <label htmlFor="email-input" className={styles.label}>
          Enter your email and we will send you a link to reset your password
          <input
            id="email-input"
            type="email"
            placeholder="you@example.com"
            className={classnames('form-control', styles['email-input'])}
            value={email}
            onChange={e => {
              const val = e.target.value;

              setEmail(val);
            }}
          />
        </label>
      </div>

      <Button
        type="submit"
        size="normal"
        kind="first"
        stretch
        className={authStyles['auth-button']}
        isBusy={submitting}
      >
        Send password reset email
      </Button>
    </form>
  );
};

export default PasswordResetRequestForm;
