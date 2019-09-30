import React from 'react';

import styles from './Placeholder.scss';

interface Props {}

const Placeholder: React.SFC<Props> = () => {
  return <div className={styles.wrapper} />;
};

export default Placeholder;
