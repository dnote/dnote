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

import PageToolbar from '../../Common/PageToolbar';
import SortMenu from './SortMenu';
import { Sort } from '../types';
import styles from './Toolbar.scss';

interface Props {
  digestUUID: string;
  sort: Sort;
  isFetched: boolean;
}

const Toolbar: React.FunctionComponent<Props> = ({
  digestUUID,
  sort,
  isFetched
}) => {
  return (
    <PageToolbar wrapperClassName={styles.wrapper}>
      <SortMenu digestUUID={digestUUID} sort={sort} disabled={!isFetched} />
    </PageToolbar>
  );
};

export default Toolbar;
