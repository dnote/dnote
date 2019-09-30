import React from 'react';
import classnames from 'classnames';

import Paginator from './Paginator';
import styles from './Top.scss';

type Position = 'top' | 'bottom';

interface Props {
  position?: Position;
}

const TopActions: React.SFC<Props> = ({ position }) => {
  return (
    <div
      className={classnames(styles.wrapper, {
        [styles.bottom]: position === 'bottom'
      })}
    >
      <Paginator />
    </div>
  );
};

export default TopActions;
