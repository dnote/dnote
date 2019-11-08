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

import React, { useState } from 'react';
import classnames from 'classnames';

import styles from './Toggle.scss';

interface Props {
  checked: boolean;
  onChange: (boolean) => void;
  label: React.ReactNode;
  id?: string;
  disabled?: boolean;
}

const Toggle: React.FunctionComponent<Props> = ({
  id,
  checked,
  onChange,
  disabled,
  label
}) => {
  return (
    <div>
      <label
        htmlFor={id}
        className={classnames(styles.label, {
          [styles.enabled]: checked,
          [styles.disabled]: !checked
        })}
      >
        <input
          id={id}
          type="checkbox"
          checked={checked}
          onChange={e => {
            onChange(e.target.checked);
          }}
          disabled={disabled}
        />

        <div className={classnames(styles.toggle, {})}>
          <div className={styles.indicator}></div>
        </div>

        {label}
      </label>
    </div>
  );
};

export default Toggle;
