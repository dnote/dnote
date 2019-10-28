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

import { countries } from 'web/libs/countries';

import styles from './CountrySelect.module.scss';

interface Props {
  id: string;
  className: string;
  onChange: (string) => void;
  value: string;
}

const CountrySelect: React.SFC<Props> = ({
  id,
  className,
  onChange,
  value
}) => {
  return (
    <div className={styles.wrapper}>
      <select
        id={id}
        className={classnames(className, styles.select, 'form-select')}
        value={value}
        onChange={onChange}
      >
        <option value="" />

        {countries.map(country => {
          return (
            <option key={country.code} value={country.code}>
              {country.name}
            </option>
          );
        })}
      </select>
    </div>
  );
};

export default CountrySelect;
