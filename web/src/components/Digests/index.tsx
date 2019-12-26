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
