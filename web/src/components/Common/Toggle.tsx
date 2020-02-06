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

import classnames from 'classnames';
import React from 'react';
import styles from './Toggle.scss';

interface Props {
  checked: boolean;
  onChange: (boolean) => void;
  label?: React.ReactNode;
  id?: string;
  disabled?: boolean;
  wrapperClassName?: string;
  kind: string;
}

export enum ToggleKind {
  first = 'first',
  green = 'green'
}

const Toggle: React.FunctionComponent<Props> = ({
  id,
  checked,
  onChange,
  disabled,
  label,
  wrapperClassName,
  kind
}) => {
  return (
    <div className={wrapperClassName}>
      <label
        htmlFor={id}
        className={classnames(styles.label, {
          [styles.first]: kind === ToggleKind.first,
          [styles.green]: kind === ToggleKind.green,
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
          <div className={styles.indicator} />
        </div>

        {label}
      </label>
    </div>
  );
};

export default Toggle;
