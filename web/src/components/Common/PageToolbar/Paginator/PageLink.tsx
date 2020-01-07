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
import classnames from 'classnames';
import { Location } from 'history';
import { Link } from 'react-router-dom';

import CaretIcon from '../../../Icons/Caret';
import { useFilters } from '../../../../store';
import styles from './Paginator.scss';

type Direction = 'next' | 'prev';

interface Props {
  direction: Direction;
  disabled: boolean;
  getPath: (page: number) => Location;
  className?: string;
}

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

const PageLink: React.FunctionComponent<Props> = ({
  direction,
  getPath,
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
      to={getPath(page)}
      aria-label={label}
      className={classnames(styles.link, className)}
    >
      {renderCaret(direction, 'black')}
    </Link>
  );
};

export default PageLink;
