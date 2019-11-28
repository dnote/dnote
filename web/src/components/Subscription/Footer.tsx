import React from 'react';

import styles from './Subscription.scss';

interface Props {}

const Footer: React.FunctionComponent<Props> = () => {
  return (
    <footer className={styles.footer}>
      &copy; 2019 Monomax Software Pty Ltd
    </footer>
  );
};

export default Footer;
