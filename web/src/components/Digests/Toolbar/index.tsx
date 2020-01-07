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

import { getDigestsPath } from 'web/libs/paths';
import PageToolbar from '../../Common/PageToolbar';
import Paginator from '../../Common/PageToolbar/Paginator';
import StatusMenu from './StatusMenu';
import { Status } from '../types';
import styles from './Toolbar.scss';

interface Props {
  total: number;
  page: number;
  status: Status;
}

const PER_PAGE = 30;

const Toolbar: React.FunctionComponent<Props> = ({ total, page, status }) => {
  return (
    <PageToolbar wrapperClassName={styles.toolbar}>
      <StatusMenu status={status} />

      <Paginator
        perPage={PER_PAGE}
        total={total}
        currentPage={page}
        getPath={(p: number) => {
          return getDigestsPath({ page: p });
        }}
      />
    </PageToolbar>
  );
};

export default Toolbar;
