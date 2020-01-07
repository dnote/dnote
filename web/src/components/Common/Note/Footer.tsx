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

import { NoteData } from 'jslib/operations/types';
import Time from '../../Common/Time';
import { nanosecToMillisec } from '../../../helpers/time';
import formatTime from '../../../helpers/time/format';
import { timeAgo } from '../../../helpers/time';
import styles from './Note.scss';

function formatAddedOn(ts: number): string {
  const ms = nanosecToMillisec(ts);
  const d = new Date(ms);

  return formatTime(d, '%MMMM %DD, %YYYY');
}

interface Props {
  note: NoteData;
  useTimeAgo: boolean;
  collapsed?: boolean;
  actions?: React.ReactElement;
}

const Footer: React.FunctionComponent<Props> = ({
  collapsed,
  actions,
  note,
  useTimeAgo
}) => {
  if (collapsed) {
    return null;
  }

  let timeText;
  if (useTimeAgo) {
    timeText = timeAgo(nanosecToMillisec(note.addedOn));
  } else {
    timeText = formatAddedOn(note.addedOn);
  }

  return (
    <footer className={styles.footer}>
      <Time
        id="note-ts"
        text={timeText}
        ms={nanosecToMillisec(note.addedOn)}
        wrapperClassName={styles.ts}
        tooltipAlignment="left"
        tooltipDirection="bottom"
      />

      {actions}
    </footer>
  );
};

export default Footer;
