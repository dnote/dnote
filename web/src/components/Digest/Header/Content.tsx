import React from 'react';

import { pluralize } from 'web/libs/string';
import { DigestData, DigestNoteData } from 'jslib/operations/types';
import Time from '../../Common/Time';
import formatTime from '../../../helpers/time/format';
import { getDigestTitle } from '../helpers';
import Progress from './Progress';
import styles from './Content.scss';

function formatCreatedAt(d: Date) {
  const now = new Date();

  const currentYear = now.getFullYear();
  const year = d.getFullYear();

  if (currentYear === year) {
    return formatTime(d, '%MMM %DD');
  }

  return formatTime(d, '%MMM %DD, %YYYY');
}

function getViewedCount(notes: DigestNoteData[]): number {
  let count = 0;

  for (let i = 0; i < notes.length; ++i) {
    const n = notes[i];

    if (n.isReviewed) {
      count++;
    }
  }

  return count;
}

interface Props {
  digest: DigestData;
}

const Content: React.FunctionComponent<Props> = ({ digest }) => {
  const viewedCount = getViewedCount(digest.notes);

  return (
    <div className={styles.header}>
      <div>
        <h1 className="page-heading">{getDigestTitle(digest)}</h1>
        <div className={styles.meta}>
          Contains {pluralize('note', digest.notes.length, true)}
          <span className={styles.sep}>&middot;</span>
          Created on{' '}
          <Time
            id="digest-ts"
            text={formatCreatedAt(new Date(digest.createdAt))}
            ms={new Date(digest.createdAt).getTime()}
            tooltipAlignment="left"
            tooltipDirection="bottom"
          />
        </div>
      </div>

      <Progress total={digest.notes.length} current={viewedCount} />
    </div>
  );
};

export default Content;
