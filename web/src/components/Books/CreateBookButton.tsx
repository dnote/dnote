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

import React, { Fragment, useState, useEffect, useRef } from 'react';
import classnames from 'classnames';
import BookPlusIcon from '../Icons/BookPlus';
import styles from './CreateBookButton.scss';

interface Props {
  disabled: boolean;
  openModal: () => void;
  className?: string;
  id?: string;
}

const CreateBookButton: React.SFC<Props> = ({
  id,
  disabled,
  openModal,
  className
}) => {
  return (
    <button
      id={id}
      type="button"
      className={classnames(
        'button-no-ui button-link',
        styles['create-button'],
        className
      )}
      disabled={disabled}
      onClick={() => {
        openModal();
      }}
    >
      <span className={styles['create-button-content']}>
        <BookPlusIcon id={`${id}-icon`} width={16} height={16} fill="#6f53c0" />
        <span className={styles['create-button-text']}>Create book</span>
      </span>
    </button>
  );
};

export default React.memo(CreateBookButton);
