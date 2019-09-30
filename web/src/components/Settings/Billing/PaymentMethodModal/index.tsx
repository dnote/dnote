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
import { StripeProvider, Elements } from 'react-stripe-elements';

import Modal, { Header, Body } from '../../../Common/Modal';
import Flash from '../../../Common/Flash';
import Form from './Form';

interface Props {
  isOpen: boolean;
  onDismiss: () => void;
  setSuccessMsg: (string) => void;
  stripe: any;
}

const PaymentMethodModal: React.SFC<Props> = ({
  isOpen,
  onDismiss,
  setSuccessMsg,
  stripe
}) => {
  const [nameOnCard, setNameOnCard] = useState('');
  const [billingCountry, setBillingCountry] = useState('');
  const [inProgress, setInProgress] = useState(false);
  const [errMessage, setErrMessage] = useState('');

  const labelId = 'payment-method-modal';

  function handleDismiss() {
    setNameOnCard('');
    setBillingCountry('');
    onDismiss();
  }

  return (
    <Modal isOpen={isOpen} onDismiss={handleDismiss} ariaLabelledBy={labelId}>
      <Header
        labelId={labelId}
        heading="Update payment method"
        onDismiss={onDismiss}
      />

      <Flash
        when={errMessage !== ''}
        kind="danger"
        onDismiss={() => {
          setErrMessage('');
        }}
      >
        {errMessage}
      </Flash>

      <Body>
        <StripeProvider stripe={stripe}>
          <Elements>
            <Form
              nameOnCard={nameOnCard}
              setNameOnCard={setNameOnCard}
              billingCountry={billingCountry}
              setBillingCountry={setBillingCountry}
              inProgress={inProgress}
              onDismiss={handleDismiss}
              setSuccessMsg={setSuccessMsg}
              setInProgress={setInProgress}
              setErrMessage={setErrMessage}
            />
          </Elements>
        </StripeProvider>
      </Body>
    </Modal>
  );
};

export default PaymentMethodModal;
