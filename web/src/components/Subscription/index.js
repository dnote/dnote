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
import { connect } from 'react-redux';
import Helmet from 'react-helmet';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import Spinner from '../Icons/Spinner';
import Plan from './Plan';
import Button from '../Common/Button';
import ServerIcon from '../Icons/Server';
import GlobeIcon from '../Icons/Globe';
import BoxIcon from '../Icons/Box';
import Flash from '../Common/Flash';

import * as paymentService from '../../services/payment';
import { getMe } from '../../services/users';
import { homePath } from '../../libs/paths';
import { updateMessage } from '../../actions/ui';
import { receiveUser } from '../../actions/auth';
import { useScript } from '../../libs/hooks';

import styles from './Subscription.module.scss';

const selfHostedPerks = [
  {
    id: 'own-machine',
    icon: <BoxIcon width="16" height="16" fill="#6e6e6e" />,
    value: 'Host on your own machine'
  }
];

const proPerks = [
  {
    id: 'hosted',
    icon: <ServerIcon width="16" height="16" fill="#245fc5" />,
    value: 'Fully hosted and managed'
  },
  {
    id: 'support',
    icon: <GlobeIcon width="16" height="16" fill="#245fc5" />,
    value: 'Support the Dnote community and development'
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

function Subscription({ userData, history, doUpdateMessage, doReceiveUser }) {
  const [openingCheckout, setOpeningCheckout] = useState(false);
  const [transacting, setTransacting] = useState(false);
  const [stripeLoaded, stripeLoadError] = useScript(
    'https://checkout.stripe.com/checkout.js'
  );

  function hideOverlay() {
    document.body.classList.remove('no-scroll');
    setTransacting(false);
  }

  const user = userData.data;

  function handlePayment() {
    if (user.cloud) {
      return;
    }

    let key;
    if (__PRODUCTION__) {
      key = 'pk_live_xvouPZFPDDBSIyMUSLZwkXfR';
    } else {
      key = 'pk_test_5926f65DQoIilZeNOiKydfoN';
    }

    setOpeningCheckout(true);

    const handler = StripeCheckout.configure({
      key,
      image: 'https://s3.amazonaws.com/dnote-asset/images/logo-circle.png',
      locale: 'auto',
      token: async token => {
        try {
          await paymentService.createSubscription({ token });
        } catch (err) {
          hideOverlay();
          console.log('Payment error', err);
          alert('error happened with payment', err);

          return;
        }

        try {
          const u = await getMe();
          doReceiveUser(u);
          doUpdateMessage('Welcome to Dnote Cloud!', 'info');
          history.push(homePath({}, { demo: false }));
        } catch (err) {
          // gracefully handle error by simply redirecting to home
          console.log('error getting user', err.message);
          window.location = '/';
        }
      },
      opened: () => {
        document.body.classList.add('no-scroll');

        setTransacting(true);
        setOpeningCheckout(true);
      },
      closed: () => {
        hideOverlay();
        setOpeningCheckout(false);
      }
    });

    handler.open({
      name: 'Dnote Pro',
      description: 'An encrypted home for your knowledge',
      amount: 300,
      currency: 'usd',
      panelLabel: '{{amount}} monthly'
    });
  }

  function renderPlanCta() {
    if (user && user.cloud) {
      return (
        <Link to="/" className="button button-third-outline button-stretch">
          Go to your notes
        </Link>
      );
    }

    return (
      <Button
        id="T-unlock-pro-btn"
        type="button"
        onClick={handlePayment}
        className={classnames('button button-third button-stretch', {
          busy: openingCheckout
        })}
        disabled={openingCheckout}
        isBusy={openingCheckout || !stripeLoaded}
      >
        Unlock
      </Button>
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

      {transacting && (
        <div className={styles['transaction-overlay']}>
          <Spinner width={50} height={50} />
        </div>
      )}

      {stripeLoadError && (
        <Flash type="danger">Stripe failed to load {stripeLoadError}</Flash>
      )}

      <div className={styles.hero}>
        <div className="container">
          <h1 className={styles.heading}>
            You can self-host or sign up for the hosted version.
          </h1>
        </div>
      </div>

      <div className="container">
        <div className={styles['plans-wrapper']}>
          <Plan
            name="Core"
            price="Free"
            features={baseFeatures}
            perks={selfHostedPerks}
            ctaContent={
              <a
                href="https://github.com/dnote/dnote"
                target="_blank"
                rel="noopener noreferrer"
                className="button button-second-outline button-stretch"
              >
                See source code
              </a>
            }
            wrapperClassName={styles['core-plan']}
          />

          <Plan
            name="Pro"
            price="$3"
            interval="month"
            features={proFeatures}
            perks={proPerks}
            ctaContent={renderPlanCta()}
            wrapperClassName={styles['pro-plan']}
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

const mapDispatchToProps = {
  doReceiveUser: receiveUser,
  doUpdateMessage: updateMessage
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Subscription);
