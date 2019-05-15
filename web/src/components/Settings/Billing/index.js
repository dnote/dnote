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
import { connect } from 'react-redux';

import Header from '../../Common/Page/Header';
import Body from '../../Common/Page/Body';
import Flash from '../../Common/Flash';
import CancelPlanModal from './CancelPlanModal';
import {
  getSubscription,
  clearSubscription,
  getSource,
  clearSource
} from '../../../actions/auth';
import SettingRow from '../SettingRow';
import PlanRow from './PlanRow';
import Placeholder from './Placeholder';
import * as paymentService from '../../../services/payment';

import settingsStyles from '../Settings.module.scss';

function ReactivateRow({
  subscriptionId,
  setSuccessMsg,
  setFailureMsg,
  doGetSubscription
}) {
  const [inProgress, setInProgress] = useState(false);

  return (
    <SettingRow
      name="Reactivate your plan"
      desc="You can reactivate your plan if you have changed your mind."
      actionContent={
        <button
          className={classnames('button-no-ui', settingsStyles.edit)}
          type="button"
          disabled={inProgress}
          onClick={() => {
            setInProgress(true);

            paymentService
              .reactivateSubscription({ subscriptionId })
              .then(() => {
                return doGetSubscription().then(() => {
                  setSuccessMsg(
                    'Your plan was reactivated. The billing cycle will be the same.'
                  );
                });
              })
              .catch(err => {
                setFailureMsg(
                  `Failed to reactivate the plan. Error: ${
                    err.message
                  }. Please contact sung@dnote.io.`
                );
                setInProgress(false);
              });
          }}
        >
          {inProgress ? 'Reactivating...' : 'Reactivate plan'}
        </button>
      }
    />
  );
}

function CancelRow({ setIsPlanModalOpen }) {
  return (
    <SettingRow
      name="Cancel current plan"
      desc="If you cancel, the plan will expire at the end of current billing period."
      actionContent={
        <button
          className={classnames('button-no-ui', settingsStyles.edit)}
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

function PaymentMethodRow({ source }) {
  let value;
  if (source.brand) {
    value = `${source.brand} ending in ${source.last4}. expiry ${
      source.exp_month
    }/${source.exp_year}`;
  } else {
    value = 'No payment method';
  }

  return <SettingRow name="Payment method" value={value} />;
}

function Content({
  subscription,
  source,
  setIsPlanModalOpen,
  successMsg,
  failureMsg,
  setSuccessMsg,
  setFailureMsg,
  doGetSubscription
}) {
  return (
    <div className="container-wide">
      <div className="row">
        <div className="col-12 col-md-12 col-lg-10">
          {successMsg && (
            <Flash
              type="success"
              wrapperClassName={settingsStyles.flash}
              onDismiss={() => {
                setSuccessMsg('');
              }}
            >
              {successMsg}
            </Flash>
          )}
          {failureMsg && (
            <Flash
              type="danger"
              wrapperClassName={settingsStyles.flash}
              onDismiss={() => {
                setFailureMsg('');
              }}
            >
              {failureMsg}
            </Flash>
          )}
        </div>

        <div className="col-12 col-md-12 col-lg-10">
          <section className={settingsStyles.section}>
            <h2 className={settingsStyles['section-heading']}>Plan</h2>

            <PlanRow subscription={subscription} />

            {subscription.id && !subscription.cancel_at_period_end && (
              <CancelRow setIsPlanModalOpen={setIsPlanModalOpen} />
            )}
            {subscription.id && subscription.cancel_at_period_end && (
              <ReactivateRow
                subscriptionId={subscription.id}
                setSuccessMsg={setSuccessMsg}
                setFailureMsg={setFailureMsg}
                doGetSubscription={doGetSubscription}
              />
            )}
          </section>

          <section className={settingsStyles.section}>
            <h2 className={settingsStyles['section-heading']}>Payment</h2>

            <PaymentMethodRow source={source} />
          </section>
        </div>
      </div>
    </div>
  );
}

function Billing({
  doGetSubscription,
  doClearSubscription,
  subscriptionData,
  sourceData,
  doGetSource,
  doClearSource
}) {
  const [isPlanModalOpen, setIsPlanModalOpen] = useState(false);
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');

  useEffect(() => {
    doGetSubscription();
    doGetSource();

    return () => {
      doClearSubscription();
      doClearSource();
    };
  }, [doGetSubscription, doClearSubscription, doGetSource, doClearSource]);

  const subscription = subscriptionData.data;
  const source = sourceData.data;

  return (
    <div>
      <Helmet>
        <title>Billing</title>
      </Helmet>

      <Header heading="Billing" />

      <Body>
        {subscriptionData.errorMessage && (
          <div className="container-wide">
            <div className="row">
              <div className="col-12 col-md-12 col-lg-10">
                <Flash type="danger" wrapperClassName={settingsStyles.flash}>
                  <div>Failed to fetch the billing information</div>
                  {subscriptionData.errorMessage}
                </Flash>
              </div>
            </div>
          </div>
        )}
        {sourceData.errorMessage && (
          <div className="container-wide">
            <div className="row">
              <div className="col-12 col-md-12 col-lg-10">
                <Flash type="danger" wrapperClassName={settingsStyles.flash}>
                  <div>Failed to fetch the payment source</div>
                  {sourceData.errorMessage}
                </Flash>
              </div>
            </div>
          </div>
        )}

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
            doGetSubscription={doGetSubscription}
          />
        )}
      </Body>

      <CancelPlanModal
        isOpen={isPlanModalOpen}
        onDismiss={() => {
          setIsPlanModalOpen(false);
        }}
        subscriptionId={subscription.id}
        setSuccessMsg={setSuccessMsg}
        setFailureMsg={setFailureMsg}
        doGetSubscription={doGetSubscription}
      />
    </div>
  );
}

function mapStateToProps(state) {
  return {
    subscriptionData: state.auth.subscription,
    sourceData: state.auth.source
  };
}

const mapDispatchToProps = {
  doGetSubscription: getSubscription,
  doClearSubscription: clearSubscription,
  doGetSource: getSource,
  doClearSource: clearSource
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Billing);
