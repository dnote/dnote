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

import React, { useState, useEffect } from 'react';
import Helmet from 'react-helmet';
import classnames from 'classnames';

import { useScript } from 'web/libs/hooks';
import { useSelector, useDispatch } from '../../../store';
import Flash from '../../Common/Flash';
import CancelPlanModal from './CancelPlanModal';
import PaymentMethodModal from './PaymentMethodModal';
import {
  getSubscription,
  clearSubscription,
  getSource,
  clearSource
} from '../../../store/auth';
import SettingRow from '../SettingRow';
import ReactivateRow from './ReactivateRow';
import PlanRow from './PlanRow';
import Placeholder from './Placeholder';
import styles from '../Settings.scss';

function CancelRow({ setIsPlanModalOpen }) {
  return (
    <SettingRow
      name="Cancel current plan"
      desc="If you cancel, the plan will expire at the end of current billing period."
      actionContent={
        <button
          className={classnames('button-no-ui', styles.edit)}
          type="button"
          onClick={() => {
            setIsPlanModalOpen(true);
          }}
        >
          Cancel plan
        </button>
      }
    />
  );
}

function PaymentMethodRow({
  stripeLoaded,
  source,
  setIsPaymentMethodModalOpen
}) {
  let value;
  if (source.brand) {
    value = `${source.brand} ending in ${source.last4}. expiry ${source.exp_month}/${source.exp_year}`;
  } else {
    value = 'No payment method';
  }

  return (
    <SettingRow
      id="T-payment-method-row"
      name="Payment method"
      value={value}
      actionContent={
        <button
          id="T-update-payment-method-button"
          className={classnames('button-no-ui', styles.edit)}
          type="button"
          onClick={() => {
            setIsPaymentMethodModalOpen(true);
          }}
          disabled={!stripeLoaded}
        >
          Update
        </button>
      }
    />
  );
}

interface ContentProps {
  subscription: any;
  source: any;
  setIsPlanModalOpen: (boolean) => void;
  setIsPaymentMethodModalOpen: (boolean) => void;
  successMsg: string;
  failureMsg: string;
  setSuccessMsg: (string) => void;
  setFailureMsg: (string) => void;
  stripeLoaded: boolean;
}

const Content: React.SFC<ContentProps> = ({
  subscription,
  source,
  setIsPlanModalOpen,
  setIsPaymentMethodModalOpen,
  successMsg,
  failureMsg,
  setSuccessMsg,
  setFailureMsg,
  stripeLoaded
}) => {
  return (
    <div>
      <Flash
        when={successMsg !== ''}
        kind="success"
        wrapperClassName={styles.flash}
        onDismiss={() => {
          setSuccessMsg('');
        }}
      >
        {successMsg}
      </Flash>
      <Flash
        when={failureMsg !== ''}
        kind="danger"
        wrapperClassName={styles.flash}
        onDismiss={() => {
          setFailureMsg('');
        }}
      >
        {failureMsg}
      </Flash>

      <div className={styles.wrapper}>
        <section className={styles.section}>
          <h2 className={styles['section-heading']}>Plan</h2>

          <PlanRow subscription={subscription} />

          {subscription.id && !subscription.cancel_at_period_end && (
            <CancelRow setIsPlanModalOpen={setIsPlanModalOpen} />
          )}
          {subscription.id && subscription.cancel_at_period_end && (
            <ReactivateRow
              subscriptionId={subscription.id}
              setSuccessMsg={setSuccessMsg}
              setFailureMsg={setFailureMsg}
            />
          )}
        </section>

        <section className={styles.section}>
          <h2 className={styles['section-heading']}>Payment</h2>

          <PaymentMethodRow
            source={source}
            setIsPaymentMethodModalOpen={setIsPaymentMethodModalOpen}
            stripeLoaded={stripeLoaded}
          />
        </section>
      </div>
    </div>
  );
};

const Billing: React.SFC = () => {
  const [isPlanModalOpen, setIsPlanModalOpen] = useState(false);
  const [isPaymentMethodModalOpen, setIsPaymentMethodModalOpen] = useState(
    false
  );
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');
  const [stripeLoaded, stripeLoadError] = useScript('https://js.stripe.com/v3');
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(getSubscription());
    dispatch(getSource());

    return () => {
      dispatch(clearSubscription());
      dispatch(clearSource());
    };
  }, [dispatch]);

  const { subscriptionData, sourceData } = useSelector(state => {
    return {
      subscriptionData: state.auth.subscription,
      sourceData: state.auth.source
    };
  });

  const subscription = subscriptionData.data;
  const source = sourceData.data;

  const key = `${__STRIPE_PUBLIC_KEY__}`;

  let stripe = null;
  if (stripeLoaded) {
    stripe = (window as any).Stripe(key);
  }

  return (
    <div>
      <Helmet>
        <title>Billing</title>
      </Helmet>

      <Flash
        when={subscriptionData.errorMessage !== ''}
        kind="danger"
        wrapperClassName={styles.flash}
      >
        <div>Failed to fetch the billing information</div>
        {subscriptionData.errorMessage}
      </Flash>

      <Flash
        when={sourceData.errorMessage !== ''}
        kind="danger"
        wrapperClassName={styles.flash}
      >
        <div>Failed to fetch the payment source</div>
        {sourceData.errorMessage}
      </Flash>

      <Flash
        when={stripeLoadError !== ''}
        kind="danger"
        wrapperClassName={styles.flash}
      >
        <div>Failed to load Stripe</div>
        {stripeLoadError}
      </Flash>

      {!subscriptionData.isFetched || !sourceData.isFetched ? (
        <Placeholder />
      ) : (
        <Content
          subscription={subscription}
          source={source}
          setIsPlanModalOpen={setIsPlanModalOpen}
          successMsg={successMsg}
          failureMsg={failureMsg}
          setSuccessMsg={setSuccessMsg}
          setFailureMsg={setFailureMsg}
          setIsPaymentMethodModalOpen={setIsPaymentMethodModalOpen}
          stripeLoaded={stripeLoaded}
        />
      )}

      <CancelPlanModal
        isOpen={isPlanModalOpen}
        onDismiss={() => {
          setIsPlanModalOpen(false);
        }}
        subscriptionId={subscription.id}
        setSuccessMsg={setSuccessMsg}
        setFailureMsg={setFailureMsg}
      />

      <PaymentMethodModal
        isOpen={isPaymentMethodModalOpen}
        onDismiss={() => {
          setIsPaymentMethodModalOpen(false);
        }}
        setSuccessMsg={setSuccessMsg}
        stripe={stripe}
      />
    </div>
  );
};

// function mapStateToProps(state) {
//   return {
//     subscriptionData: state.auth.subscription,
//     sourceData: state.auth.source
//   };
// }

// const mapDispatchToProps = {
//   doGetSubscription: getSubscription,
//   doClearSubscription: clearSubscription,
//   doGetSource: getSource,
//   doClearSource: clearSource
// };

export default Billing;
