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

import React, { useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import Helmet from 'react-helmet';

import { usePrevious } from 'web/libs/hooks';
import { parseSearchString } from 'jslib/helpers/url';
import { useDispatch, useSelector } from '../../store';
import { getDigests } from '../../store/digests';
import { Status } from './types';
import Flash from '../Common/Flash';
import List from './List';
import Toolbar from './Toolbar';

function useFetchDigests(params: { page: number; status: Status }) {
  const dispatch = useDispatch();

  const prevParams = usePrevious(params);

  useEffect(() => {
    if (
      !prevParams ||
      prevParams.page !== params.page || prevParams.status !== params.status
    ) {
      dispatch(getDigests(params));
    }
  }, [dispatch, params, prevParams]);
}

interface Props extends RouteComponentProps {}

const Digests: React.FunctionComponent<Props> = ({ location }) => {
  const { digests } = useSelector(state => {
    return {
      digests: state.digests
    };
  });
  const { page, status } = parseSearchString(location.search);
  useFetchDigests({
    page: page || 1,
    status
  });

  return (
    <div className="page page-mobile-full">
      <Helmet>
        <title>Digests</title>
      </Helmet>

      <div className="container mobile-fw">
        <div className="page-header">
          <h1 className="page-heading">Digests</h1>
        </div>

        <Flash kind="danger" when={Boolean(digests.errorMessage)}>
          Error getting digests: {digests.errorMessage}
        </Flash>
      </div>

      <div className="container mobile-nopadding">
        <Toolbar total={digests.total} page={digests.page} status={status} />

        <List
          isFetching={digests.isFetching}
          isFetched={digests.isFetched}
          items={digests.data}
        />
      </div>
    </div>
  );
};

export default Digests;
