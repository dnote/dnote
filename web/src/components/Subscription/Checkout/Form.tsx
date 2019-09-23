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
import { withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';
import Helmet from 'react-helmet';
import { injectStripe } from 'react-stripe-elements';
import { History } from 'history';

import services from 'web/libs/services';
import { getHomePath } from 'web/libs/paths';
import Sidebar from './Sidebar';
import Flash from '../../Common/Flash';
import { getCurrentUser } from '../../../store/auth';
import { useDispatch } from '../../../store';
import { setMessage } from '../../../store/ui';
import NameOnCardInput from '../../Common/PaymentInput/NameOnCard';
import CardInput from '../../Common/PaymentInput/Card';
import CountryInput from '../../Common/PaymentInput/Country';
import styles from './Form.scss';

interface Props extends RouteComponentProps {
  stripe: any;
  stripeLoadError: string;
  history: History;
}

const Form: React.SFC<Props> = ({ stripe, stripeLoadError, history }) => {
  const [nameOnCard, setNameOnCard] = useState('');
  const cardElementRef = useRef(null);
  const [cardElementLoaded, setCardElementLoaded] = useState(false);
  const [billingCountry, setBillingCountry] = useState('');
  const [transacting, setTransacting] = useState(false);
  const [errMessage, setErrMessage] = useState('');
  const dispatch = useDispatch();

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

      await services.payment.createSubscription({
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

    await dispatch(getCurrentUser());

    const homePath = getHomePath();

    dispatch(
      setMessage({
        message: 'Welcome to Dnote Pro',
        kind: 'info',
        path: homePath.pathname
      })
    );
    history.push(homePath);
  }

  return (
    <div className={styles.wrapper}>
      <form className={classnames('container')} onSubmit={handleSubmit}>
        <Helmet>
          <title>Subscriptions</title>
        </Helmet>

        <Flash
          when={errMessage !== ''}
          kind="danger"
          wrapperClassName={styles.flash}
          onDismiss={() => {
            setErrMessage('');
          }}
        >
          {errMessage}
        </Flash>
        <Flash
          when={stripeLoadError !== ''}
          kind="danger"
          wrapperClassName={styles.flash}
        >
          Failed to load stripe. {stripeLoadError}
        </Flash>

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
    </div>
  );
};

export default injectStripe(withRouter(Form));
