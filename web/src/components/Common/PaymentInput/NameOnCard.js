import React from 'react';

import classnames from 'classnames';

import styles from './PaymentInput.module.scss';

function NameOnCard({ value, onUpdate, containerClassName, labelClassName }) {
  return (
    <div className={classnames(containerClassName)}>
      <label htmlFor="name-on-card" className="label-full">
        <span className={classnames(labelClassName)}>Name on Card</span>
        <input
          autoFocus
          id="name-on-card"
          className={classnames(
            'text-input text-input-stretch text-input-medium',
            styles.input
          )}
          type="text"
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
