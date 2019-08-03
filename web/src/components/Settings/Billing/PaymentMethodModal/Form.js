import React, { useState, useRef } from 'react';
import { injectStripe } from 'react-stripe-elements';

import Button from '../../../Common/Button';
import NameOnCardInput from '../../../Common/PaymentInput/NameOnCard';
import CardInput from '../../../Common/PaymentInput/Card';
import CountryInput from '../../../Common/PaymentInput/Country';

import settingsStyles from '../../Settings.module.scss';
import * as paymentService from '../../../../services/payment';
import styles from './Form.module.scss';

function Form({
  stripe,
  nameOnCard,
  setNameOnCard,
  billingCountry,
  setBillingCountry,
  inProgress,
  onDismiss,
  setSuccessMsg,
  setInProgress,
  doGetSource,
  setErrMessage
}) {
  const [cardElementLoaded, setCardElementLoaded] = useState(false);
  const cardElementRef = useRef(null);

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

      await paymentService.updateSource({ source, country: billingCountry });
      await doGetSource();

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
          isBusy={!cardElementLoaded || inProgress}
        >
          Update
        </Button>

        <Button
          type="button"
          kind="second"
          disabled={inProgress}
          onClick={onDismiss}
        >
          Cancel
        </Button>
      </div>
    </form>
  );
}

export default injectStripe(Form);
