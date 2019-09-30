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

import React from 'react';
import classnames from 'classnames';

import CloseIcon from '../../Icons/Close';
import styles from './ModalHeader.scss';

interface Props {
  labelId: string;
  heading: string;
  onDismiss: () => void;
}

const Header: React.SFC<Props> = ({ labelId, heading, onDismiss }) => {
  return (
    <div className={styles.wrapper}>
      <strong id={labelId}>{heading}</strong>

      <button
        onClick={onDismiss}
        type="button"
        aria-label="Close the modal"
        className={classnames('button-no-ui T-modal-close', styles.button)}
      >
        <CloseIcon width={16} height={16} />
      </button>
    </div>
  );
};

export default Header;
