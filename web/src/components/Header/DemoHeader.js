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
import classnames from 'classnames';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import { getSubscriptionPath, getJoinPath } from 'web/libs/paths';

import styles from './DemoHeader.module.scss';

function getPricingPath(user) {
  if (user) {
    return getSubscriptionPath();
  }

  return getJoinPath({ referrer: getSubscriptionPath() });
}

function DemoHeader({ user }) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.left}>
        <div className={styles.heading}>Live Demo</div>

        <span className={styles.subheading}>
          <span className={styles.support}>
            Get your own encrypted repository of knowledge.
          </span>
        </span>
      </div>

      <div className={styles.right}>
        <Link
          to={getPricingPath(user)}
          className={classnames(
            styles.cta,
            'button button-normal button-second'
          )}
        >
          Get started
        </Link>
        <a
          href="/"
          className={classnames(
            styles['quit-mobile'],
            'button button-normal button-first'
          )}
        >
          Quit demo
        </a>
      </div>

      <a href="/" className={styles.quit}>
        Quit demo
      </a>
    </div>
  );
}

function mapStateToProps(state) {
  return {
    user: state.auth.user
  };
}

export default connect(mapStateToProps)(DemoHeader);
