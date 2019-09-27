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
import Helmet from 'react-helmet';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import { getSubscriptionCheckoutPath } from 'web/libs/paths';
import ProPlan from './Plan/Pro';
import CorePlan from './Plan/Core';
import FeatureList from './FeatureList';
import { useSelector } from '../../store';

import styles from './Subscription.scss';

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

interface Props {}

const Subscription: React.SFC<Props> = () => {
  const { user } = useSelector(state => {
    return {
      user: state.auth.user.data
    };
  });

  function renderPlanCta() {
    if (user && user.pro) {
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
};

export default Subscription;
