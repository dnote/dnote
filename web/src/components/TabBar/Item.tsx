import React from 'react';

import styles from './Item.scss';

interface Props {}

const Item: React.SFC<Props> = ({ children }) => {
  return <li className={styles.wrapper}>{children}</li>;
};

export default Item;
