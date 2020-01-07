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
import Helmet from 'react-helmet';

import { NoteState } from '../../store/note';
import { nanosecToMillisec } from '../../helpers/time';
import formatTime from '../../helpers/time/format';

function formatAddedOn(ts: number): string {
  const ms = nanosecToMillisec(ts);
  const d = new Date(ms);

  return formatTime(d, '%MMM %D %YYYY');
}

function getTitle(note: NoteState): string {
  if (!note.isFetched) {
    return 'Note';
  }

  return `Note: ${note.data.book.label} (${formatAddedOn(note.data.addedOn)})`;
}

interface Props {
  note: NoteState;
}

const HeaderData: React.FunctionComponent<Props> = ({ note }) => {
  const title = getTitle(note);
  return (
    <Helmet>
      <title>{title}</title>
    </Helmet>
  );
};

export default HeaderData;
