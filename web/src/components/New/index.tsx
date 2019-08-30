import React, { useState } from 'react';
import { Prompt, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';
import Helmet from 'react-helmet';
import { withRouter } from 'react-router-dom';

import Editor from '../Common/Editor';
import Flash from '../Common/Flash';
import { useDispatch, useSelector } from '../../store';
import { resetEditor } from '../../store/editor';
import { createBook } from '../../store/books';
import { setMessage } from '../../store/ui';
import * as notesOperation from '../../operations/notes';
import { getNotePath, notePathDef } from '../../libs/paths';
import { useCleanupEditor, useFocusTextarea } from '../../libs/hooks/editor';
import styles from './New.scss';

interface Props extends RouteComponentProps {}

const New: React.SFC<Props> = ({ history }) => {
  const { editor } = useSelector(state => {
    return {
      editor: state.editor
    };
  });
  const dispatch = useDispatch();

  const [errMessage, setErrMessage] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [bookSelectorOpen, setBookSelectorOpen] = useState(!editor.bookUUID);
  const [textareaEl, setTextareaEl] = useState(null);

  useCleanupEditor();
  useFocusTextarea(textareaEl, bookSelectorOpen);

  return (
    <div
      className={classnames(
        styles.container,
        'container mobile-nopadding page'
      )}
    >
      <Helmet>
        <title>New</title>
      </Helmet>

      <Flash kind="danger" when={Boolean(errMessage)}>
        Error: {errMessage}
      </Flash>

      <div className={styles.wrapper}>
        <div
          className={classnames(styles.overlay, {
            [styles.active]: bookSelectorOpen
          })}
        />
        <div className={styles.header}>
          <h2 className={styles.heading}>New note</h2>
        </div>

        <Editor
          isNew
          isBusy={submitting}
          setBookSelectorOpen={setBookSelectorOpen}
          bookSelectorOpen={bookSelectorOpen}
          textareaEl={textareaEl}
          setTextareaEl={setTextareaEl}
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

              const res = await notesOperation.create({
                bookUUID,
                content: draftContent
              });

              dispatch(resetEditor());

              const dest = getNotePath(res.result.uuid);
              history.push(dest);

              dispatch(
                setMessage({
                  message: 'Created a note',
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

      <Prompt
        message="You have unsaved changes. Continue?"
        when={editor.dirty}
      />
    </div>
  );
};

export default React.memo(withRouter(New));
