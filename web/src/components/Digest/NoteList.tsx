import React from 'react';
import Helmet from 'react-helmet';

import { DigestData } from 'jslib/operations/types';
import { DigestNoteData } from 'jslib/operations/types';
import { getDigestTitle } from './helpers';
import { useDispatch } from '../../store';
import { setDigestNoteReviewed } from '../../store/digest';
import Placeholder from '../Common/Note/Placeholder';
import NoteItem from './NoteItem';
import styles from './Digest.scss';

interface Props {
  notes: DigestNoteData[];
  digest: DigestData;
  isFetched: boolean;
  isFetching: boolean;
}

const NoteList: React.FunctionComponent<Props> = ({
  isFetched,
  isFetching,
  notes,
  digest
}) => {
  const dispatch = useDispatch();

  function handleSetReviewed(noteUUID: string, isReviewed: boolean) {
    return dispatch(
      setDigestNoteReviewed({ digestUUID: digest.uuid, noteUUID, isReviewed })
    );
  }

  if (isFetching) {
    return (
      <div className={styles.wrapper}>
        <Placeholder wrapperClassName={styles.item} />
        <Placeholder wrapperClassName={styles.item} />
        <Placeholder wrapperClassName={styles.item} />
      </div>
    );
  }
  if (!isFetched) {
    return null;
  }

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>{`${getDigestTitle(digest)} - Digest`}</title>
      </Helmet>
      <ul className={styles.list}>
        {notes.map(note => {
          return (
            <NoteItem
              key={note.uuid}
              note={note}
              onSetReviewed={handleSetReviewed}
            />
          );
        })}
      </ul>
    </div>
  );
};

export default NoteList;
