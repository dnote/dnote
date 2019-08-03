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
import { injectStripe } from 'react-stripe-elements';

import Sidebar from './Sidebar';
import Flash from '../../Common/Flash';
import * as paymentService from '../../../services/payment';
import { getCurrentUser } from '../../../actions/auth';
import { updateMessage } from '../../../actions/ui';
import { getHomePath } from '../../../libs/paths';
import NameOnCardInput from '../../Common/PaymentInput/NameOnCard';
import CardInput from '../../Common/PaymentInput/Card';
import CountryInput from '../../Common/PaymentInput/Country';

import styles from './Form.module.scss';

function Form({
  stripe,
  stripeLoadError,
  doGetCurrentUser,
  doUpdateMessage,
  history
}) {
  const [nameOnCard, setNameOnCard] = useState('');
  const cardElementRef = useRef(null);
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
              <NameOnCardInput
                value={nameOnCard}
                onUpdate={setNameOnCard}
                containerClassName={styles['input-row']}
                labelClassName={styles.label}
              />

              <CardInput
                cardElementRef={cardElementRef}
                setCardElementLoaded={setCardElementLoaded}
                containerClassName={styles['input-row']}
                labelClassName={styles.label}
              />

              <CountryInput
                value={billingCountry}
                onUpdate={setBillingCountry}
                containerClassName={styles['input-row']}
                labelClassName={styles.label}
              />
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
