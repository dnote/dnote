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

import styles from './PaymentInput.scss';

interface Props {
  value: string;
  onUpdate: (string) => void;
  containerClassName?: string;
  labelClassName?: string;
}

const NameOnCard: React.SFC<Props> = ({
  value,
  onUpdate,
  containerClassName,
  labelClassName
}) => {
  return (
    <div className={classnames(containerClassName)}>
      <label htmlFor="name-on-card" className="label-full">
        <span className={classnames(labelClassName)}>Name on Card</span>
        <input
          autoFocus
          id="name-on-card"
          className={classnames(
            'text-input text-input-stretch text-input-medium',
            styles.input
          )}
          type="text"
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

export default NameOnCard;
