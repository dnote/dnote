/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

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

const stickyThresholdY = 24;

function checkSticky(y: number): boolean {
  return y > stickyThresholdY;
}

const Header: React.FunctionComponent<Props> = ({ digest, isFetched }) => {
  const [isSticky, setIsSticky] = useState(false);

  function handleScroll() {
    const y = getScrollYPos();
    const nextSticky = checkSticky(y);

    if (nextSticky) {
      setIsSticky(true);
    } else if (!nextSticky) {
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
      <div className="container mobile-fw">
        <div className={styles.header}>
          {isFetched ? <Content digest={digest} /> : <Placeholder />}
        </div>
      </div>
    </div>
  );
};

export default Header;
