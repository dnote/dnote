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

import Button from '../../Common/Button';
import ServerIcon from '../../Icons/Server';
import GlobeIcon from '../../Icons/Globe';

import styles from './Sidebar.scss';

const perks = [
  {
    id: 'hosted',
    icon: <ServerIcon width={16} height={16} fill="#4d4d8b" />,
    value: 'Fully hosted and managed'
  },
  {
    id: 'support',
    icon: <GlobeIcon width={16} height={16} fill="#4d4d8b" />,
    value: 'Support the Dnote community and development'
  }
];

function Sidebar({ isReady, transacting }) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.header}>
        <div className={styles['plan-name']}>Pro</div>

        <ul className={classnames('list-unstyled', styles.perks)}>
          {perks.map(perk => {
            return (
              <li key={perk.id} className={styles['perk-item']}>
                <div className={styles['perk-icon']}>{perk.icon}</div>
                <div className={styles['perk-value']}>{perk.value}</div>
              </li>
            );
          })}
        </ul>

        <div className={styles['price-wrapper']}>
          <strong className={styles.price}>$3.00</strong>
          <div className={styles.interval}>/ month</div>
        </div>

        <Button
          id="T-purchase-button"
          type="submit"
          kind="first"
          size="normal"
          className={classnames(
            'button button-large button-third button-stretch',
            styles['purchase-button']
          )}
          disabled={transacting}
          isBusy={transacting || !isReady}
        >
          Purchase Dnote Pro
        </Button>
      </div>

      <p className={styles.assurance}>You can cancel your plan any time.</p>
    </div>
  );
}

export default Sidebar;
