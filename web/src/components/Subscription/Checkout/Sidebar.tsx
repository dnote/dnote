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
import classnames from 'classnames';

import Button from '../../Common/Button';
import Toggle, { ToggleKind } from '../../Common/Toggle';
import PaymentSummary from './PaymentSummary';
import Price from './Price';
import styles from './Sidebar.scss';

interface Props {
  isReady: boolean;
  transacting: boolean;
  yearly: boolean;
  setYearly: (boolean) => null;
}

function Sidebar({ isReady, transacting, yearly, setYearly }) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.header}>
        <div className={styles['plan-name']}>Dnote Pro</div>
        <p className={styles['plan-desc']}>Fully hosted and managed</p>

        <Price yearly={yearly} />

        <div className={styles['interval-toggle-wrapper']}>
          <button
            type="button"
            onClick={() => {
              setYearly(false);
            }}
            className={classnames(
              'button-no-ui button-no-padding',
              styles['interval-label'],
              {
                [styles.active]: !yearly
              }
            )}
          >
            Bill monthly
          </button>
          <Toggle
            id="T-yearly-toggle"
            checked={yearly}
            onChange={() => {
              setYearly(!yearly);
            }}
            disabled={transacting}
            wrapperClassName={styles['interval-toggle']}
            kind={ToggleKind.first}
          />
          <button
            type="button"
            onClick={() => {
              setYearly(true);
            }}
            className={classnames(
              'button-no-ui button-no-padding',
              styles['interval-label'],
              {
                [styles.active]: yearly
              }
            )}
          >
            Bill yearly
          </button>
        </div>

        <PaymentSummary yearly={yearly} />

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

      <p className={styles.assurance}>You can cancel auto-renewal any time.</p>
    </div>
  );
}

export default Sidebar;
