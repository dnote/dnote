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

import ProPlan from './Plan/Pro';
import CorePlan from './Plan/Core';
import FeatureList from './FeatureList';
import Footer from './Footer';
import { useSelector } from '../../store';

import styles from './Subscription.scss';

const proFeatures = [
  {
    id: 'spaced-rep',
    label: <div>Automated Spaced Repetition</div>
  },
  {
    id: 'email-support',
    label: <div>Email support</div>
  }
];

const coreFeatures = [
  {
    id: 'oss',
    label: <div>Open source</div>
  },
  {
    id: 'num-notes',
    label: <div>Unlimited notes</div>
  },
  {
    id: 'num-books',
    label: <div>Unlimited books</div>
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
    id: 'web',
    label: <div>Web application</div>
  },
  {
    id: 'ext',
    label: <div>Chrome/Firefox extension</div>
  },
  {
    id: 'forum-support',
    label: <div>Forum support</div>
  }
];

interface Props {}

const Subscription: React.FunctionComponent<Props> = () => {
  const { user } = useSelector(state => {
    return {
      user: state.auth.user.data
    };
  });

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>Subscriptions</title>
        <meta
          name="description"
          content="Dnote has a simple pricing with a free core and a paid addon."
        />
      </Helmet>

      <div className={styles.content}>
        <div className={styles.hero}>
          <div className="container">
            <h1 className={styles.heading}>Choose your Dnote plan.</h1>
          </div>
        </div>

        <div className="container">
          <div className={styles['plans-wrapper']}>
            <CorePlan
              wrapperClassName={styles['core-plan']}
              user={user}
              bottomContent={
                <div className={styles.bottom}>
                  <FeatureList features={coreFeatures} />
                </div>
              }
            />

            <ProPlan
              wrapperClassName={styles['pro-plan']}
              user={user}
              bottomContent={
                <div className={styles.bottom}>
                  <div className={styles['pro-prelude']}>
                    Everything from the core plan, plus:
                  </div>
                  <FeatureList
                    features={proFeatures}
                    wrapperClassName={styles['pro-feature-list']}
                  />
                </div>
              }
            />
          </div>
        </div>
      </div>

      <Footer />
    </div>
  );
};

export default Subscription;
