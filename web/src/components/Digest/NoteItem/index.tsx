import React, { Fragment, useState } from 'react';
import { Link } from 'react-router-dom';

import { DigestNoteData } from 'jslib/operations/types';
import { getNotePath } from 'web/libs/paths';
import Note from '../../Common/Note';
import Flash from '../../Common/Flash';
import NoteItemHeader from './Header';
import styles from '../Digest.scss';

interface Props {
  note: DigestNoteData;
  onSetReviewed: (string, boolean) => Promise<any>;
}

const NoteItem: React.FunctionComponent<Props> = ({ note, onSetReviewed }) => {
  const [collapsed, setCollapsed] = useState(note.isReviewed);
  const [errorMessage, setErrMessage] = useState('');

  return (
    <li className={styles.item}>
      <Note
        collapsed={collapsed}
        note={note}
        header={
          <Fragment>
            <NoteItemHeader
              note={note}
              collapsed={collapsed}
              setCollapsed={setCollapsed}
              onSetReviewed={onSetReviewed}
              setErrMessage={setErrMessage}
            />

            <Flash kind="danger" when={errorMessage !== ''}>
              {errorMessage}
            </Flash>
          </Fragment>
        }
        footerActions={
          <Link className={styles.action} to={getNotePath(note.uuid)}>
            Go to note ›
          </Link>
        }
        footerUseTimeAgo
      />
    </li>
  );
};

export default NoteItem;
