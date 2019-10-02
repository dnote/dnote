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

import React, { useState, useRef } from 'react';
import { Prompt, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';
import { withRouter } from 'react-router-dom';

import operations from 'web/libs/operations';
import { getEditorSessionkey } from 'web/libs/editor';
import { getNotePath, notePathDef } from 'web/libs/paths';
import { useFocusTextarea } from 'web/libs/hooks/editor';
import Editor from '../Common/Editor';
import { useDispatch, useSelector } from '../../store';
import { resetEditor, EditorSession } from '../../store/editor';
import { createBook } from '../../store/books';
import { setMessage } from '../../store/ui';
import styles from '../New/New.scss';

interface Props extends RouteComponentProps {
  noteUUID: string;
  persisted: boolean;
  editor: EditorSession;
  setErrMessage: React.Dispatch<string>;
}

const Edit: React.SFC<Props> = ({
  noteUUID,
  persisted,
  editor,
  history,
  setErrMessage
}) => {
  const { prevLocation } = useSelector(state => {
    return {
      prevLocation: state.route.prevLocation
    };
  });
  const dispatch = useDispatch();
  const [submitting, setSubmitting] = useState(false);
  const textareaRef = useRef(null);

  useFocusTextarea(textareaRef.current);

  return (
    <div className={styles.wrapper}>
      <div className={classnames(styles.overlay, {})} />
      <div className={styles.header}>
        <h2 className={styles.heading}>Edit note</h2>
      </div>

      <Editor
        editor={editor}
        isBusy={submitting}
        textareaRef={textareaRef}
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

            dispatch(resetEditor(editor.sessionKey));

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

      <Prompt message="You have unsaved changes. Continue?" when={!persisted} />
    </div>
  );
};

export default React.memo(withRouter(Edit));
