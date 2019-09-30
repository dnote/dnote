import React from 'react';
import classnames from 'classnames';

import CloseIcon from '../../Icons/Close';
import CaretIcon from '../../Icons/CaretSolid';
import styles from './SearchInput.scss';

interface Props {
  onReset?: (Event) => void;
  resetShown?: boolean;
  expanded?: boolean;
  setExpanded?: (boolean) => void;
}

const Actions: React.SFC<Props> = ({
  onReset,
  resetShown,
  setExpanded,
  expanded
}) => {
  const resettable = Boolean(onReset);
  const expandable = Boolean(setExpanded);

  return (
    <div className={styles.actions}>
      {resettable && (
        <button
          type="button"
          className={classnames('button-no-ui', styles.reset, {
            [styles['reset-shown']]: resetShown
          })}
          aria-label="Clear search"
          onClick={onReset}
        >
          <CloseIcon width={16} height={16} />
        </button>
      )}
      {expandable && (
        <button
          type="button"
          className={classnames('button-no-ui', styles.expand)}
          aria-label="Expand search menu"
          onClick={() => {
            setExpanded(!expanded);
          }}
        >
          <CaretIcon width={12} height={12} />
        </button>
      )}
    </div>
  );
};

export default Actions;
