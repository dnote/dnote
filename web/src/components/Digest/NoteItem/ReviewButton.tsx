import React, { useState } from 'react';
import classnames from 'classnames';

import digestStyles from '../Digest.scss';
import styles from './ReviewButton.scss';

interface Props {
  noteUUID: string;
  isReviewed: boolean;
  setCollapsed: (boolean) => void;
  onSetReviewed: (string, boolean) => Promise<any>;
  setErrMessage: (string) => void;
}

const ReviewButton: React.FunctionComponent<Props> = ({
  noteUUID,
  isReviewed,
  setCollapsed,
  onSetReviewed,
  setErrMessage
}) => {
  const [checked, setChecked] = useState(isReviewed);

  return (
    <label className={styles.wrapper}>
      <input
        type="checkbox"
        checked={checked}
        onChange={e => {
          const val = e.target.checked;

          // update UI optimistically
          setErrMessage('');
          setChecked(val);
          setCollapsed(val);

          onSetReviewed(noteUUID, val).catch(err => {
            // roll back the UI update in case of error
            setChecked(!val);
            setCollapsed(!val);

            setErrMessage(err.message);
          });
        }}
      />
      <span className={classnames(digestStyles['header-action'], styles.text)}>
        Reviewed
      </span>
    </label>
  );
};

export default ReviewButton;
