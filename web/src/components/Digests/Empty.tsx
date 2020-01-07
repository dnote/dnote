import React from 'react';
import { Link } from 'react-router-dom';

import { getRepetitionsPath } from 'web/libs/paths';
import styles from './Empty.scss';

interface Props {}

const Empty: React.FunctionComponent<Props> = () => {
  return (
    <div className={styles.wrapper}>
      <h3>No digests were found.</h3>

      <p className={styles.support}>
        You could <Link to={getRepetitionsPath()}>create repetition rules</Link>{' '}
        first.
      </p>

      <p className={styles['md-support']}>
        Digests are automatically created based on your repetition rules.
      </p>
    </div>
  );
};

export default Empty;
