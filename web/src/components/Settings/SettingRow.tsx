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

import styles from './SettingRow.scss';

interface Props {
  name: string;
  actionContent: React.ReactNode;
  id?: string;
  desc?: string;
  value?: string;
}

const SettingRow: React.SFC<Props> = ({
  name,
  desc,
  value,
  actionContent,
  id
}) => {
  return (
    <div className={classnames(styles.wrapper, styles.row)} id={id}>
      <div>
        <h3 className={styles.name}>{name}</h3>
        <p className={styles.desc}>{desc}</p>
      </div>

      <div className={styles.right}>
        {value}
        <div className={styles.action}>{actionContent}</div>
      </div>
    </div>
  );
};

export default SettingRow;
