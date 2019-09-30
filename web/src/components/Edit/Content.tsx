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
  const textareaRef = useRef(null);

  useFocusTextarea(textareaRef.current);
  useCleanupEditor();

  return (
    <div className={styles.wrapper}>
      <div className={classnames(styles.overlay, {})} />
      <div className={styles.header}>
        <h2 className={styles.heading}>Edit note</h2>
      </div>

      <Editor
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
