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
import { Location } from 'history';

import PageLink from './PageLink';
import styles from './Paginator.scss';

interface Props {
  perPage: number;
  total: number;
  currentPage: number;
  getPath: (page: number) => Location;
}

function getMaxPage(total: number, perPage: number): number {
  if (total === 0) {
    return 1;
  }

  return Math.ceil(total / perPage);
}

const Paginator: React.FunctionComponent<Props> = ({
  perPage,
  total,
  currentPage,
  getPath
}) => {
  const hasNext = currentPage * perPage < total;
  const hasPrev = currentPage > 1;
  const maxPage = getMaxPage(total, perPage);

  return (
    <nav className={styles.wrapper}>
      <span className={styles.info}>
        <span className={styles.label}>{currentPage}</span> of{' '}
        <span className={styles.label}>{maxPage}</span>
      </span>

      <PageLink
        direction="prev"
        disabled={!hasPrev}
        className={styles['link-prev']}
        getPath={getPath}
      />
      <PageLink direction="next" disabled={!hasNext} getPath={getPath} />
    </nav>
  );
};

export default Paginator;
