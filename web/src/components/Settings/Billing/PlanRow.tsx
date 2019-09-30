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
import { Link } from 'react-router-dom';
import moment from 'moment';
import classnames from 'classnames';

import { getPlanLabel } from 'web/libs/subscription';
import LogoIcon from '../../Icons/Logo';
import styles from './PlanRow.scss';
import settingRowStyles from '../SettingRow.scss';

function getPlanPeriodMessage(subscription: any): string {
  if (!subscription.id) {
    return 'You do not have a subscription.';
  }

  const label = getPlanLabel(subscription);

  const endDate = moment.unix(subscription.current_period_end);

  if (subscription.cancel_at_period_end) {
    return `Your ${label} plan will end on ${endDate.format(
      'YYYY MMM Do'
    )} and will not renew.`;
  }

  const renewDate = endDate.add(1, 'day');
  return `Your ${label} plan will renew on ${renewDate.format('YYYY MMM Do')}.`;
}

interface Props {
  subscription: any;
}

const PlanRow: React.SFC<Props> = ({ subscription }) => {
  return (
    <div className={classnames(settingRowStyles.row, styles.wrapper)}>
      <div className={styles.content}>
        <LogoIcon width={40} height={40} />
        <div className={styles.detail}>
          <div>
            <strong className={styles.label}>
              {getPlanLabel(subscription)}
            </strong>

            {subscription.cancel_at_period_end && (
              <span className={styles.status}>(cancelled)</span>
            )}
          </div>
          <p className={styles.desc}>{getPlanPeriodMessage(subscription)}</p>
        </div>
      </div>

      <div>
        {!subscription.id && (
          <Link
            className="button button-normal button-first"
            to="/subscriptions"
          >
            Upgrade
          </Link>
        )}
      </div>
    </div>
  );
};

export default PlanRow;
