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

import Item from './Item';
import { getNewPath, getBooksPath, getRepetitionsPath } from 'web/libs/paths';
import { Filters, toSearchObj } from 'jslib/helpers/filters';
import styles from './Nav.scss';

interface Props {
  filters: Filters;
}

const Nav: React.SFC<Props> = ({ filters }) => {
  const searchObj = toSearchObj(filters);

  return (
    <nav className={styles.wrapper}>
      <ul className={classnames('list-unstyled', styles.list)}>
        <Item to={getNewPath(searchObj)} label="New" />
        <Item to={getBooksPath(searchObj)} label="Books" />
        <Item to={getRepetitionsPath(searchObj)} label="Repetition" />
        {/* <Item to={getRandomPath(searchObj)} label="Random" /> */}
      </ul>
    </nav>
  );
};

export default Nav;
