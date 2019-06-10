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

import React, { useState, useRef } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import classnames from 'classnames';
import Helmet from 'react-helmet';
import { injectStripe, CardElement } from 'react-stripe-elements';

import Sidebar from './Sidebar';
import CountrySelect from './CountrySelect';
import Flash from '../../Common/Flash';
import * as paymentService from '../../../services/payment';
import { getCurrentUser } from '../../../actions/auth';
import { updateMessage } from '../../../actions/ui';
import { getHomePath } from '../../../libs/paths';

import styles from './Form.module.scss';

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

function Form({
  stripe,
  stripeLoadError,
  doGetCurrentUser,
  doUpdateMessage,
  history
}) {
  const [nameOnCard, setNameOnCard] = useState('');
  const cardElementRef = useRef(null);
  const [cardElementFocused, setCardElementFocused] = useState(false);
  const [cardElementLoaded, setCardElementLoaded] = useState(false);
  const [billingCountry, setBillingCountry] = useState('');
  const [transacting, setTransacting] = useState(false);
  const [errMessage, setErrMessage] = useState('');

  async function handleSubmit(e) {
    e.preventDefault();

    if (!cardElementLoaded) {
      return;
    }
    if (!nameOnCard) {
      setErrMessage('Please enter the name on card');
      return;
    }
    if (!billingCountry) {
      setErrMessage('Please enter the country');
      return;
    }

    setTransacting(true);

    try {
      const { source, error } = await stripe.createSource({
        type: 'card',
        currency: 'usd',
        owner: {
          name: nameOnCard
        }
      });

      if (error) {
        throw error;
      }

      await paymentService.createSubscription({
        source,
        country: billingCountry
      });
    } catch (err) {
      console.log('error subscribing', err);
      setTransacting(false);
      setErrMessage(err.message);
      return;
    }

    setNameOnCard('');
    setBillingCountry('');
    cardElementRef.current.clear();
    setTransacting(false);

    await doGetCurrentUser();
    doUpdateMessage('Welcome to Dnote Pro', 'info');
    history.push(getHomePath({}, { demo: false }));
  }

  return (
    <form
      className={classnames('container', styles.wrapper)}
      onSubmit={handleSubmit}
    >
      <Helmet>
        <title>Subscriptions</title>
      </Helmet>

      {errMessage && (
        <Flash
          type="danger"
          wrapperClassName={styles.flash}
          onDismiss={() => {
            setErrMessage('');
          }}
        >
          {errMessage}
        </Flash>
      )}
      {stripeLoadError && (
        <Flash type="danger" wrapperClassName={styles.flash}>
          Failed to load stripe. {stripeLoadError}
        </Flash>
      )}

      <div className="row">
        <div className="col-12 col-lg-7 col-xl-8">
          <div className={styles['content-wrapper']}>
            <h1 className={styles.heading}>You are almost there.</h1>

            <div className={styles.content}>
              <div className={styles['input-row']}>
                <label htmlFor="name-on-card" className="label-full">
                  <span className={styles.label}>Name on Card</span>
                  <input
                    autoFocus
                    id="name-on-card"
                    className={classnames(
                      'text-input text-input-stretch text-input-medium',
                      styles.input
                    )}
                    type="text"
                    value={nameOnCard}
                    onChange={e => {
                      const val = e.target.value;
                      setNameOnCard(val);
                    }}
                  />
                </label>
              </div>

              <div
                className={classnames(styles['input-row'], styles['card-row'])}
              >
                {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
                <label htmlFor="card-number" className={styles.number}>
                  <span className={styles.label}>Card Number</span>

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
                      cardElementRef.current = el;
                      setCardElementLoaded(true);
                    }}
                    style={elementStyles}
                  />
                </label>
              </div>

              <div className={styles['input-row']}>
                {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
                <label htmlFor="billing-country" className="label-full">
                  <span className={styles.label}>Country</span>
                  <CountrySelect
                    id="billing-country"
                    className={classnames(
                      styles['countries-select'],
                      styles.input
                    )}
                    value={billingCountry}
                    onChange={e => {
                      const val = e.target.value;
                      setBillingCountry(val);
                    }}
                  />
                </label>
              </div>
            </div>
          </div>
        </div>

        <div className="col-12 col-lg-5 col-xl-4">
          <Sidebar isReady={cardElementLoaded} transacting={transacting} />
        </div>
      </div>
    </form>
  );
}

const mapDispatchToProps = {
  doGetCurrentUser: getCurrentUser,
  doUpdateMessage: updateMessage
};

export default injectStripe(
  withRouter(
    connect(
      null,
      mapDispatchToProps
    )(Form)
  )
);
