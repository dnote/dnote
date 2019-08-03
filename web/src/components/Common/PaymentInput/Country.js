import React from 'react';

import classnames from 'classnames';

import CountrySelect from './CountrySelect';
import styles from './PaymentInput.module.scss';

function NameOnCard({ value, onUpdate, containerClassName, labelClassName }) {
  return (
    <div className={classnames(containerClassName)}>
      {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
      <label htmlFor="billing-country" className="label-full">
        <span className={classnames(labelClassName)}>Country</span>
        <CountrySelect
          id="billing-country"
          className={classnames(styles['countries-select'], styles.input)}
          value={value}
          onChange={e => {
            const val = e.target.value;
            onUpdate(val);
          }}
        />
      </label>
    </div>
  );
}

export default NameOnCard;
