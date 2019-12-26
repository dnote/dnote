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
