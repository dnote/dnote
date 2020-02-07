/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import React from 'react';
import { StripeProvider, Elements } from 'react-stripe-elements';
// import { Redirect } from 'react-router-dom';

import { useScript } from 'web/libs/hooks';
import CheckoutForm from './Form';
// import { getHomePath } from '../../../libs/paths';
// import { useSelector } from '../../../store/hooks';

const Checkout: React.FunctionComponent = () => {
  const [stripeLoaded, stripeLoadError] = useScript('https://js.stripe.com/v3');
  const [paypalLoaded, paypalLoadError] = useScript(
    `https://www.paypal.com/sdk/js?client-id=${__PAYPAL_CLIENT_ID__}&vault=true`
  );

  console.log(paypalLoaded);

  const key = `${__STRIPE_PUBLIC_KEY__}`;

  let stripe = null;
  if (stripeLoaded) {
    stripe = (window as any).Stripe(key);
  }

  // TODO: reactivate this guard clause after grandfathered free users aree migrated
  // const { user, userFetched } = useSelector(state => {
  //   return {
  //     user: state.auth.user.data,
  //     userFetched: state.auth.user.isFetched
  //   };
  // });
  // if (userFetched && user.pro) {
  //   return <Redirect to={getHomePath()} />;
  // }

  return (
    <StripeProvider stripe={stripe}>
      <Elements>
        <CheckoutForm
          stripeLoadError={stripeLoadError}
          paypalLoadError={paypalLoadError}
        />
      </Elements>
    </StripeProvider>
  );
};

export default Checkout;
