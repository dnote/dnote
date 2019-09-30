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

import { Link } from 'react-router-dom';
import { useFilters, useSelector } from '../../../store';
import { getHomePath } from 'web/libs/paths';
import CaretIcon from '../../Icons/Caret';
import styles from './Paginator.scss';

// PER_PAGE is the number of results per page in the response from the backend implementation's API.
// Currently it is fixed.
const PER_PAGE = 30;

type Direction = 'next' | 'prev';

const renderCaret = (direction: Direction, fill: string) => {
  return (
    <CaretIcon
      fill={fill}
      width={12}
      height={12}
      className={styles[`caret-${direction}`]}
    />
  );
};

interface PageLinkProps {
  direction: Direction;
  disabled: boolean;
  className?: string;
}

const PageLink: React.SFC<PageLinkProps> = ({
  direction,
  disabled,
  className
}) => {
  const filters = useFilters();

  if (disabled) {
    return (
      <span className={classnames(styles.link, styles.disabled, className)}>
        {renderCaret(direction, 'gray')}
      </span>
    );
  }

  let page;
  if (direction === 'next') {
    page = filters.page + 1;
  } else {
    page = filters.page - 1;
  }

  let label;
  if (direction === 'next') {
    label = 'Next page';
  } else {
    label = 'Previous page';
  }

  return (
    <Link
      to={getHomePath({
        ...filters.queries,
        page
      })}
      aria-label={label}
      className={classnames(styles.link, className)}
    >
      {renderCaret(direction, 'black')}
    </Link>
  );
};

interface PaginatorProps {}

const Paginator: React.SFC<PaginatorProps> = () => {
  const filters = useFilters();
  const { notes } = useSelector(state => {
    return {
      notes: state.notes
    };
  });

  const hasNext = filters.page * PER_PAGE < notes.total;
  const hasPrev = filters.page > 1;
  const maxPage = Math.ceil(notes.total / PER_PAGE);

  let currentPage;
  if (maxPage > 0) {
    currentPage = filters.page;
  } else {
    currentPage = 0;
  }

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
      />
      <PageLink direction="next" disabled={!hasNext} />
    </nav>
  );
};

export default Paginator;
