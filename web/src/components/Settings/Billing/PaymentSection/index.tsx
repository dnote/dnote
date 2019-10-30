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

import React, { Fragment } from 'react';
import classnames from 'classnames';

import PaymentMethodRow from './PaymentMethodRow';
import settingsStyles from '../../Settings.scss';
import { SourceData } from '../../../../store/auth/type';
import Placeholder from './Placeholder';
import styles from './Placeholder.scss';

interface Props {
  source: SourceData;
  setIsPaymentMethodModalOpen: (boolean) => void;
  stripeLoaded: boolean;
  isFetched: boolean;
}

const PaymentSection: React.FunctionComponent<Props> = ({
  source,
  setIsPaymentMethodModalOpen,
  stripeLoaded,
  isFetched
}) => {
  if (!isFetched) {
    return <Placeholder />;
  }

  return (
    <Fragment>
      <PaymentMethodRow
        source={source}
        setIsPaymentMethodModalOpen={setIsPaymentMethodModalOpen}
        stripeLoaded={stripeLoaded}
      />
    </Fragment>
  );
};

export default PaymentSection;
