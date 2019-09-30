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
import Helmet from 'react-helmet';

import { NoteState } from '../../store/note';
import { nanosecToMillisec, getShortMonthName } from '../../helpers/time';

function formatAddedOn(ts: number): string {
  const ms = nanosecToMillisec(ts);
  const d = new Date(ms);

  const month = getShortMonthName(d);
  const date = d.getDate();
  const year = d.getFullYear();

  return `${month} ${date} ${year}`;
}

function getTitle(note: NoteState): string {
  if (!note.isFetched) {
    return 'Note';
  }

  return `Note (${formatAddedOn(note.data.added_on)}) in ${
    note.data.book.label
  }`;
}

function getDescription(note: NoteState): string {
  if (!note.isFetched) {
    return 'View microlessons and write your own.';
  }

  const book = note.data.book;
  return `View microlessons in ${book.label} and write your own. Dnote is a home for your everyday learning.`;
}

interface Props {
  note: NoteState;
}

const HeaderData: React.SFC<Props> = ({ note }) => {
  const title = getTitle(note);
  const description = getDescription(note);

  const noteData = note.data;

  return (
    <Helmet>
      <title>{title}</title>
      <meta name="description" content={description} />
      <meta name="twitter:card" content="summary" />
      <meta name="twitter:title" content={title} />
      <meta name="twitter:description" content={noteData.content} />
      <meta
        name="twitter:image"
        content="https://s3.amazonaws.com/dnote-assets/images/bf3fed4fb122e394e26bcf35d63e26f8.png"
      />
      <meta
        name="og:image"
        content="https://s3.amazonaws.com/dnote-assets/images/bf3fed4fb122e394e26bcf35d63e26f8.png"
      />
      <meta name="og:title" content={title} />
      <meta name="og:description" content={noteData.content} />
    </Helmet>
  );
};

export default HeaderData;
