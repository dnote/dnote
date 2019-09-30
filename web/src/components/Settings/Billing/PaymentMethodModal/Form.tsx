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
import { injectStripe } from 'react-stripe-elements';

import services from 'web/libs/services';
import { useDispatch } from '../../../../store';
import { getSource } from '../../../../store/auth';
import Button from '../../../Common/Button';
import NameOnCardInput from '../../../Common/PaymentInput/NameOnCard';
import CardInput from '../../../Common/PaymentInput/Card';
import CountryInput from '../../../Common/PaymentInput/Country';
import settingsStyles from '../../Settings.scss';
import styles from './Form.scss';

interface Props {
  stripe: any;
  nameOnCard: string;
  setNameOnCard: (string) => void;
  billingCountry: string;
  setBillingCountry: (string) => void;
  inProgress: boolean;
  onDismiss: () => void;
  setSuccessMsg: (string) => void;
  setInProgress: (boolean) => void;
  setErrMessage: (string) => void;
}

const Form: React.SFC<Props> = ({
  stripe,
  nameOnCard,
  setNameOnCard,
  billingCountry,
  setBillingCountry,
  inProgress,
  onDismiss,
  setSuccessMsg,
  setInProgress,
  setErrMessage
}) => {
  const [cardElementLoaded, setCardElementLoaded] = useState(false);
  const cardElementRef = useRef(null);
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

    setSuccessMsg('');
    setErrMessage('');
    setInProgress(true);

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

      await services.payment.updateSource({ source, country: billingCountry });
      await dispatch(getSource());

      setSuccessMsg('Your payment method was successfully updated.');
      setInProgress(false);
      onDismiss();
    } catch (err) {
      setErrMessage(`An error occurred: ${err.message}`);
      setInProgress(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} id="T-payment-method-form">
      <div>
        <NameOnCardInput
          value={nameOnCard}
          onUpdate={setNameOnCard}
          containerClassName={styles['input-row']}
        />

        <CardInput
          cardElementRef={cardElementRef}
          setCardElementLoaded={setCardElementLoaded}
          containerClassName={styles['input-row']}
        />

        <CountryInput
          value={billingCountry}
          onUpdate={setBillingCountry}
          containerClassName={styles['input-row']}
        />
      </div>

      <div className={settingsStyles.actions}>
        <Button
          type="submit"
          kind="first"
          size="normal"
          isBusy={!cardElementLoaded || inProgress}
        >
          Update
        </Button>

        <Button
          type="button"
          kind="second"
          size="normal"
          disabled={inProgress}
          onClick={onDismiss}
        >
          Cancel
        </Button>
      </div>
    </form>
  );
};

export default injectStripe(Form);
