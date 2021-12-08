/* Copyright (C) 2019, 2020, 2021 Monomax Software Pty Ltd
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

import classnames from 'classnames';
import React, { Fragment, useEffect, useRef, useState } from 'react';
import { Prompt, RouteComponentProps, withRouter } from 'react-router-dom';
import { focusTextarea } from 'web/libs/dom';
import { useFocus } from 'web/libs/hooks/dom';
import operations from 'web/libs/operations';
import { getNotePath, notePathDef } from 'web/libs/paths';
import { useDispatch } from '../../store';
import { createBook } from '../../store/books';
import { EditorSession, resetEditor } from '../../store/editor';
import { setMessage } from '../../store/ui';
import Editor from '../Common/Editor';
import Flash from '../Common/Flash';
import styles from './New.scss';

interface Props extends RouteComponentProps {
  editor: EditorSession;
  persisted: boolean;
}

// useInitFocus initializes the focus on HTML elements depending on the current
// state of the editor.
function useInitFocus({ bookLabel, content, textareaRef, setTriggerFocus }) {
  useEffect(() => {
    if (!bookLabel && !content) {
      setTriggerFocus();
    } else {
      const textareaEl = textareaRef.current;

      if (textareaEl && document.activeElement !== textareaEl) {
        focusTextarea(textareaEl);
      }
    }
  }, [setTriggerFocus, bookLabel, textareaRef, content]);
}

const New: React.FunctionComponent<Props> = ({
  editor,
  persisted,
  history
}) => {
  const dispatch = useDispatch();
  const [errMessage, setErrMessage] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const textareaRef = useRef(null);
  const [setTriggerFocus, triggerRef] = useFocus();

  useInitFocus({
    bookLabel: editor.bookLabel,
    content: editor.content,
    textareaRef,
    setTriggerFocus
  });

  return (
    <Fragment>
      <div
        className={classnames(
          'page page-mobile-full container mobile-nopadding',
          styles.container
        )}
      >
        <Flash kind="danger" when={Boolean(errMessage)}>
          Error: {errMessage}
        </Flash>

        <div className={styles.wrapper}>
          <div className={styles.header}>
            <h2 className={styles.heading}>New notes</h2>
          </div>

          <Editor
            isNew
            editor={editor}
            isBusy={submitting}
            textareaRef={textareaRef}
            bookSelectorTriggerRef={triggerRef}
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

                const res = await operations.notes.create({
                  bookUUID,
                  content: draftContent
                });

                dispatch(resetEditor(editor.sessionKey));

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
          when={!persisted}
        />
      </div>
    </Fragment>
  );
};

export default React.memo(withRouter(New));
