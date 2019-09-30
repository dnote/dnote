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

import React, { useState } from 'react';
import { CardElement } from 'react-stripe-elements';

import classnames from 'classnames';
import styles from './PaymentInput.scss';

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

interface Props {
  cardElementRef?: React.MutableRefObject<any>;
  setCardElementLoaded: (boolean) => void;
  containerClassName?: string;
  labelClassName?: string;
}

const Card: React.SFC<Props> = ({
  cardElementRef,
  setCardElementLoaded,
  containerClassName,
  labelClassName
}) => {
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
};

export default Card;
