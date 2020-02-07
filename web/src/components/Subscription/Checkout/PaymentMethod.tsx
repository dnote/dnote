import React from 'react';
import classnames from 'classnames';

import { PaymentMethod } from './helpers';
import styles from './Form.scss';

interface Props {
  method: PaymentMethod;
  isActive: boolean;
  setMethod: (PaymentMethod) => void;
}

const Component: React.SFC<Props> = ({
  children,
  method,
  isActive,
  setMethod
}) => {
  return (
    <label
      className={classnames(
        'button button-large button-second-outline',
        styles.method,
        {
          [styles['method-active']]: isActive
        }
      )}
    >
      <input
        type="radio"
        name="payment_method"
        onChange={() => {
          setMethod(method);
        }}
        value={method}
        checked={isActive}
      />
      {children}
    </label>
  );
};

export default Component;
