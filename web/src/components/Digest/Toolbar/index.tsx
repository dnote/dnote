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
