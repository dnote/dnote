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

import React from 'react';
import { connect } from 'react-redux';
import Helmet from 'react-helmet';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import ProPlan from './Plan/Pro';
import CorePlan from './Plan/Core';
import FeatureList from './FeatureList';
import { getSubscriptionCheckoutPath } from '../../libs/paths';

import styles from './Subscription.module.scss';

const proFeatures = [
  {
    id: 'core',
    label: <div className={styles['feature-bold']}>Everything in core</div>
  },
  {
    id: 'host',
    label: <div>Hosting</div>
  },
  {
    id: 'auto',
    label: <div>Automatic update and migration</div>
  },
  {
    id: 'email-support',
    label: <div>Email support</div>
  }
];

const baseFeatures = [
  {
    id: 'backup',
    label: <div>Encrypted backup using AES256</div>
  },
  {
    id: 'sync',
    label: <div>Multi-device sync</div>
  },
  {
    id: 'cli',
    label: <div>Command line interface</div>
  },
  {
    id: 'atom',
    label: <div>Atom plugin</div>
  },
  {
    id: 'web',
    label: <div>Web client</div>
  },
  {
    id: 'digest',
    label: <div>Automated email digest</div>
  },
  {
    id: 'ext',
    label: <div>Firefox/Chrome extension</div>
  },
  {
    id: 'foss',
    label: <div>Free and open source</div>
  },
  {
    id: 'forum-support',
    label: <div>Forum support</div>
  }
];
function Subscription({ userData }) {
  const user = userData.data;

  //  function handlePayment() {
  //    if (user.cloud) {
  //      return;
  //    }
  //
  //    let key;
  //    if (__PRODUCTION__) {
  //      key = 'pk_live_xvouPZFPDDBSIyMUSLZwkXfR';
  //    } else {
  //      key = 'pk_test_5926f65DQoIilZeNOiKydfoN';
  //    }
  //
  //    const handler = StripeCheckout.configure({
  //      key,
  //      image: 'https://s3.amazonaws.com/dnote-asset/images/logo-circle.png',
  //      locale: 'auto',
  //      token: async token => {
  //        try {
  //          await paymentService.createSubscription({ token });
  //        } catch (err) {
  //          hideOverlay();
  //          console.log('Payment error', err);
  //          alert('error happened with payment', err);
  //
  //          return;
  //        }
  //
  //        try {
  //          const u = await getMe();
  //          doReceiveUser(u);
  //          doUpdateMessage('Welcome to Dnote Cloud!', 'info');
  //          history.push(getHomePath({}, { demo: false }));
  //        } catch (err) {
  //          // gracefully handle error by simply redirecting to home
  //          console.log('error getting user', err.message);
  //          window.location = '/';
  //        }
  //      },
  //      opened: () => {
  //        document.body.classList.add('no-scroll');
  //
  //        setTransacting(true);
  //        setOpeningCheckout(true);
  //      },
  //      closed: () => {
  //        hideOverlay();
  //        setOpeningCheckout(false);
  //      }
  //    });
  //
  //    handler.open({
  //      name: 'Dnote Pro',
  //      description: 'An encrypted home for your knowledge',
  //      amount: 300,
  //      currency: 'usd',
  //      panelLabel: '{{amount}} monthly'
  //    });
  //  }

  function renderPlanCta() {
    if (user && user.cloud) {
      return (
        <Link
          to="/"
          className="button button-normal button-third-outline button-stretch"
        >
          Go to your notes
        </Link>
      );
    }

    return (
      <Link
        id="T-unlock-pro-btn"
        type="button"
        className={classnames(
          'button button-normal button-third button-stretch'
        )}
        to={getSubscriptionCheckoutPath()}
      >
        Unlock
      </Link>
    );
  }

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>Subscriptions</title>
        <meta
          name="description"
          content="Dnote has a simple pricing with a free core and a paid addon."
        />
      </Helmet>

      <div className={styles.hero}>
        <div className="container">
          <h1 className={styles.heading}>
            You can self-host or sign up for the hosted version.
          </h1>
        </div>
      </div>

      <div className="container">
        <div className={styles['plans-wrapper']}>
          <CorePlan
            wrapperClassName={styles['core-plan']}
            ctaContent={
              <a
                href="https://github.com/dnote/dnote"
                target="_blank"
                rel="noopener noreferrer"
                className="button button-normal button-second-outline button-stretch"
              >
                See source code
              </a>
            }
            bottomContent={<FeatureList features={baseFeatures} />}
          />

          <ProPlan
            wrapperClassName={styles['pro-plan']}
            ctaContent={renderPlanCta()}
            bottomContent={<FeatureList features={proFeatures} />}
          />
        </div>
      </div>
    </div>
  );
}

function mapStateToProps(state) {
  return {
    userData: state.auth.user
  };
}

export default connect(mapStateToProps)(Subscription);
