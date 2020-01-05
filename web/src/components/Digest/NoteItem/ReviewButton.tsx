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
