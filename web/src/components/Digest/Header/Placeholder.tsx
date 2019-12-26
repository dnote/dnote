import React from 'react';
import classnames from 'classnames';

import styles from './Placeholder.scss';

interface Props {}

const HeaderPlaceholder: React.FunctionComponent<Props> = () => {
  return (
    <div className={styles.wrapper}>
      <div className={classnames('holder holder-dark', styles.title)} />

      <div className={classnames('holder holder-dark', styles.meta)} />
    </div>
  );
};

export default HeaderPlaceholder;
