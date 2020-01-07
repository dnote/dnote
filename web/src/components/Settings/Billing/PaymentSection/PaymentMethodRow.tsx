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

import SettingRow from '../../SettingRow';
import { SourceData } from '../../../../store/auth/type';
import styles from '../../Settings.scss';

interface Props {
  stripeLoaded: boolean;
  source: SourceData;
  setIsPaymentMethodModalOpen: (bool) => void;
}

const PaymentMethodRow: React.FunctionComponent<Props> = ({
  stripeLoaded,
  source,
  setIsPaymentMethodModalOpen
}) => {
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
};

export default PaymentMethodRow;
