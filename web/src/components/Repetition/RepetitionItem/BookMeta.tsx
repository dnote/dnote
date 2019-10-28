import React from 'react';

import { BookDomain } from 'jslib/operations/types';
import { pluralize } from 'web/libs/string';
import styles from './RepetitionItem.scss';

interface ContentProps {
  bookDomain: BookDomain;
  bookCount: number;
}

const Content: React.FunctionComponent<ContentProps> = ({
  bookDomain,
  bookCount
}) => {
  if (bookDomain === BookDomain.All) {
    return <span>From all books</span>;
  }

  let verb;
  if (bookDomain === BookDomain.Excluding) {
    verb = 'Excluding';
  } else if (bookDomain === BookDomain.Including) {
    verb = 'From';
  }

  return (
    <span>
      {verb} {bookCount} {pluralize('book', bookCount)}
    </span>
  );
};

interface Props {
  bookDomain: BookDomain;
  bookCount: number;
}

const BookMeta: React.FunctionComponent<Props> = ({
  bookDomain,
  bookCount
}) => {
  return (
    <span className={styles['book-meta']}>
      <Content bookDomain={bookDomain} bookCount={bookCount} />
    </span>
  );
};

export default BookMeta;
