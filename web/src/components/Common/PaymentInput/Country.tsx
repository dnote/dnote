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

import CountrySelect from './CountrySelect';
import styles from './PaymentInput.scss';

interface Props {
  value: string;
  onUpdate: (string) => void;
  containerClassName?: string;
  labelClassName?: string;
}

const Country: React.SFC<Props> = ({
  value,
  onUpdate,
  containerClassName,
  labelClassName
}) => {
  return (
    <div className={classnames(containerClassName)}>
      {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
      <label htmlFor="billing-country" className="label-full">
        <span className={classnames(labelClassName)}>Country</span>
        <CountrySelect
          id="billing-country"
          className={classnames(styles['countries-select'], styles.input)}
          value={value}
          onChange={e => {
            const val = e.target.value;
            onUpdate(val);
          }}
        />
      </label>
    </div>
  );
};

export default Country;
