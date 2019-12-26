import React from 'react';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import { DigestNoteData } from 'jslib/operations/types';
import { getHomePath } from 'web/libs/paths';
import ReviewButton from './ReviewButton';
import Button from '../../Common/Button';
import CaretIcon from '../../Icons/CaretSolid';

import noteStyles from '../../Common/Note/Note.scss';
import styles from './Header.scss';

interface Props {
  note: DigestNoteData;
  setCollapsed: (boolean) => void;
  onSetReviewed: (string, boolean) => Promise<any>;
  setErrMessage: (string) => void;
  collapsed: boolean;
}

const Header: React.FunctionComponent<Props> = ({
  note,
  collapsed,
  setCollapsed,
  onSetReviewed,
  setErrMessage
}) => {
  let fill;
  if (collapsed) {
    fill = '#8c8c8c';
  } else {
    fill = '#000000';
  }

  return (
    <header className={noteStyles.header}>
      <div className={noteStyles['header-left']}>
        <Button
          className={styles['header-action']}
          type="button"
          kind="no-ui"
          onClick={() => {
            setCollapsed(!collapsed);
          }}
        >
          <CaretIcon
            fill={fill}
            width={12}
            height={12}
            className={classnames({ [styles['caret-collapsed']]: collapsed })}
          />
        </Button>

        <h1
          className={classnames(noteStyles['book-label'], styles['book-label'])}
        >
          <Link to={getHomePath({ book: note.book.label })}>
            {note.book.label}
          </Link>
        </h1>
      </div>

      <div className={noteStyles['header-right']}>
        <ReviewButton
          isReviewed={note.isReviewed}
          noteUUID={note.uuid}
          setCollapsed={setCollapsed}
          onSetReviewed={onSetReviewed}
          setErrMessage={setErrMessage}
        />
      </div>
    </header>
  );
};

export default Header;
