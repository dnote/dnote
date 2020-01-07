/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import SettingRow from '../../SettingRow';
import styles from '../../Settings.scss';

interface Props {
  setIsPlanModalOpen: (bool) => void;
}

const CancelRow: React.FunctionComponent<Props> = ({ setIsPlanModalOpen }) => {
  return (
    <SettingRow
      name="Cancel current plan"
      desc="If you cancel, the plan will expire at the end of current billing period."
      actionContent={
        <button
          className={classnames('button-no-ui', styles.edit)}
          type="button"
          onClick={() => {
            setIsPlanModalOpen(true);
          }}
        >
          Cancel plan
        </button>
      }
    />
  );
};

export default CancelRow;
