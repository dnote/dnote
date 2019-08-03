import React, { useState } from 'react';
import { CardElement } from 'react-stripe-elements';

import classnames from 'classnames';

import styles from './PaymentInput.module.scss';

const elementStyles = {
  base: {
    color: '#32325D',
    fontFamily: 'Source Code Pro, Consolas, Menlo, monospace',
    fontSize: '16px',
    fontSmoothing: 'antialiased',

    '::placeholder': {
      color: '#CFD7DF'
    },
    ':-webkit-autofill': {
      color: '#e39f48'
    }
  },
  invalid: {
    color: '#E25950',

    '::placeholder': {
      color: '#FFCCA5'
    }
  }
};

function Card({
  cardElementRef,
  setCardElementLoaded,
  containerClassName,
  labelClassName
}) {
  const [cardElementFocused, setCardElementFocused] = useState(false);

  return (
    <div className={classnames(styles['card-row'], containerClassName)}>
      {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
      <label htmlFor="card-number" className={styles.number}>
        <span className={classnames(labelClassName)}>Card Number</span>

        <CardElement
          id="card"
          className={classnames(styles['card-number'], styles.input, {
            [styles['card-number-active']]: cardElementFocused
          })}
          onFocus={() => {
            setCardElementFocused(true);
          }}
          onBlur={() => {
            setCardElementFocused(false);
          }}
          onReady={el => {
            if (cardElementRef) {
              // eslint-disable-next-line no-param-reassign
              cardElementRef.current = el;
            }
            setCardElementLoaded(true);
          }}
          style={elementStyles}
        />
      </label>
    </div>
  );
}

export default Card;
