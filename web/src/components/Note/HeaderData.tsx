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
