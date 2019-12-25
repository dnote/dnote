import React, { useEffect } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import Helmet from 'react-helmet';

import { getDigestsPath } from 'web/libs/paths';
import { usePrevious } from 'web/libs/hooks';
import { parseSearchString } from 'jslib/helpers/url';
import PageToolbar from '../Common/PageToolbar';
import { useDispatch, useSelector } from '../../store';
import { getDigests } from '../../store/digests';
import Flash from '../Common/Flash';
import List from './List';

const PER_PAGE = 30;

function useFetchDigests(page: number) {
  const dispatch = useDispatch();

  const prevPage = usePrevious(page);

  useEffect(() => {
    if (prevPage !== page) {
      dispatch(getDigests(page));
    }
  }, [dispatch, page, prevPage]);
}

interface Props extends RouteComponentProps {}

const Digests: React.SFC<Props> = ({ location }) => {
  const { digests } = useSelector(state => {
    return {
      digests: state.digests
    };
  });
  const { page } = parseSearchString(location.search);
  useFetchDigests(page || 1);

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
          Error getting notes: {digests.errorMessage}
        </Flash>
      </div>

      <div className="container mobile-nopadding">
        <PageToolbar
          perPage={PER_PAGE}
          total={digests.total}
          currentPage={digests.page}
          getPath={(p: number) => {
            return getDigestsPath({ page: p });
          }}
        />

        <List isFetched={digests.isFetched} items={digests.data} />
      </div>
    </div>
  );
};

export default Digests;
