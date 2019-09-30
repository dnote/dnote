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
import Helmet from 'react-helmet';

import { Filters } from 'jslib/helpers/filters';

interface Props {
  filters: Filters;
}

function getTitle(filters: Filters): string {
  if (filters.queries.book.length === 1) {
    return `Notes in ${filters.queries.book}`;
  }

  return 'Notes';
}

const HeaderData: React.SFC<Props> = ({ filters }) => {
  const title = getTitle(filters);

  return (
    <Helmet>
      <title>{title}</title>
    </Helmet>
  );
};

export default HeaderData;
