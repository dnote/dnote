import React from 'react';

import { SearchParams, Status } from './types';
import styles from './Empty.scss';

interface Props {
  params: SearchParams;
}

const Empty: React.FunctionComponent<Props> = ({ params }) => {
  if (params.status === Status.Unreviewed) {
    return (
      <div className={styles.wrapper}>
        You have completed reviewing this digest.
      </div>
    );
  }

  return <div className={styles.wrapper}>No results matched your filters.</div>;
};

export default Empty;
