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

import React, { useEffect } from 'react';
import { Link, RouteComponentProps } from 'react-router-dom';
import Helmet from 'react-helmet';

import { usePrevious } from 'web/libs/hooks';
import { getSubscriptionPath } from 'web/libs/paths';
import { parseSearchString } from 'jslib/helpers/url';
import { useDispatch, useSelector } from '../../store';
import { getDigests } from '../../store/digests';
import { Status } from './types';
import Flash from '../Common/Flash';
import List from './List';
import Toolbar from './Toolbar';
import styles from './Digests.scss';

function useFetchDigests(params: { page: number; status: Status }) {
  const dispatch = useDispatch();

  const prevParams = usePrevious(params);

  useEffect(() => {
    if (
      !prevParams ||
      prevParams.page !== params.page ||
      prevParams.status !== params.status
    ) {
      dispatch(getDigests(params));
    }
  }, [dispatch, params, prevParams]);
}

interface Props extends RouteComponentProps {}

const Digests: React.FunctionComponent<Props> = ({ location }) => {
  const { user, digests } = useSelector(state => {
    return {
      digests: state.digests,
      user: state.auth.user.data
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
      </div>

      <div className="container mobile-nopadding">
        <Flash
          kind="danger"
          when={Boolean(digests.errorMessage)}
          wrapperClassName={styles.flash}
        >
          Error getting digests: {digests.errorMessage}
        </Flash>

        <Flash when={!user.pro} kind="warning" wrapperClassName={styles.flash}>
          Digests are not enabled on your plan.{' '}
          <Link to={getSubscriptionPath()}>Upgrade here.</Link>
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
