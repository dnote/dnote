import React, { useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';
import { withRouter } from 'react-router-dom';

import operations from 'web/libs/operations';
import { getNotePath, notePathDef } from 'web/libs/paths';
import { useCleanupEditor, useFocusTextarea } from 'web/libs/hooks/editor';
import Editor from '../Common/Editor';
import { useDispatch, useSelector } from '../../store';
import { resetEditor } from '../../store/editor';
import { createBook } from '../../store/books';
import { setMessage } from '../../store/ui';
import styles from '../New/New.scss';

interface Props extends RouteComponentProps {
  noteUUID: string;
  setErrMessage: React.Dispatch<string>;
}

const Edit: React.SFC<Props> = ({ noteUUID, history, setErrMessage }) => {
  const { editor, prevLocation } = useSelector(state => {
    return {
      editor: state.editor,
      prevLocation: state.route.prevLocation
    };
  });
  const dispatch = useDispatch();

  const [submitting, setSubmitting] = useState(false);
  const [bookSelectorOpen, setBookSelectorOpen] = useState(!editor.bookUUID);
  const [textareaEl, setTextareaEl] = useState(null);

  useFocusTextarea(textareaEl);
  useCleanupEditor();

  return (
    <div className={styles.wrapper}>
      <div
        className={classnames(styles.overlay, {
          [styles.active]: bookSelectorOpen
        })}
      />
      <div className={styles.header}>
        <h2 className={styles.heading}>Edit note</h2>
      </div>

      <Editor
        isBusy={submitting}
        setBookSelectorOpen={setBookSelectorOpen}
        bookSelectorOpen={bookSelectorOpen}
        textareaEl={textareaEl}
        setTextareaEl={setTextareaEl}
        cancelPath={prevLocation}
        onSubmit={async ({ draftContent, draftBookUUID }) => {
          setSubmitting(true);

          try {
            let bookUUID;

            if (!draftBookUUID) {
              const book = await dispatch(createBook(editor.bookLabel));
              bookUUID = book.uuid;
            } else {
              bookUUID = draftBookUUID;
            }

            const note = await operations.notes.update(noteUUID, {
              book_uuid: bookUUID,
              content: draftContent
            });

            dispatch(resetEditor());

            const dest = getNotePath(note.uuid);
            history.push(dest);

            dispatch(
              setMessage({
                message: 'Updated the note',
                kind: 'info',
                path: notePathDef
              })
            );
          } catch (err) {
            setErrMessage(err.message);
            setSubmitting(false);
          }
        }}
      />
    </div>
  );
};

export default React.memo(withRouter(Edit));
