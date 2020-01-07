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
import { Link } from 'react-router-dom';

import { getDigestPath } from 'web/libs/paths';
import { SearchParams } from './types';
import CloseIcon from '../Icons/Close';
import styles from './ClearSearchBar.scss';

interface Props {
  params: SearchParams;
  digestUUID: string;
}

const ClearSearchBar: React.FunctionComponent<Props> = ({
  params,
  digestUUID
}) => {
  const isActive = params.sort !== '' || params.status !== '';

  if (!isActive) {
    return null;
  }

  return (
    <div className={styles.wrapper}>
      <Link className={styles.button} to={getDigestPath(digestUUID)}>
        <CloseIcon width={20} height={20} />
        <span className={styles.text}>
          Clear the current filters, and sorts
        </span>
      </Link>
    </div>
  );
};

export default ClearSearchBar;
