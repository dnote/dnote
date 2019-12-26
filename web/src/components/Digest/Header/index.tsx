import React, { useState } from 'react';
import classnames from 'classnames';

import { DigestData } from 'jslib/operations/types';
import { useEventListener } from 'web/libs/hooks';
import { getScrollYPos } from 'web/libs/dom';
import Placeholder from './Placeholder';
import Content from './Content';
import styles from './Content.scss';

interface Props {
  isFetched: boolean;
  digest: DigestData;
}

const stickyThresholdY = 0;

function checkSticky(y: number): boolean {
  return y > stickyThresholdY;
}

const Header: React.FunctionComponent<Props> = ({ digest, isFetched }) => {
  const [isSticky, setIsSticky] = useState(false);

  function handleScroll() {
    const y = getScrollYPos();

    const nextSticky = checkSticky(y);

    if (nextSticky && !isSticky) {
      setIsSticky(true);
    } else if (!nextSticky && isSticky) {
      setIsSticky(false);
    }
  }

  useEventListener(document, 'scroll', handleScroll);

  return (
    <div
      className={classnames(styles['header-container'], {
        [styles['header-sticky']]: isSticky
      })}
    >
      <div className="container">
        {isFetched ? <Content digest={digest} /> : <Placeholder />}
      </div>
    </div>
  );
};

export default Header;
