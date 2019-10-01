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

import React, { useState, useRef, useEffect, Fragment } from 'react';
import { Prompt, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';
import Helmet from 'react-helmet';
import { withRouter } from 'react-router-dom';

import { focusTextarea } from 'web/libs/dom';
import operations from 'web/libs/operations';
import { getNotePath, notePathDef } from 'web/libs/paths';
import { useCleanupEditor } from 'web/libs/hooks/editor';
import { useFocus } from 'web/libs/hooks/dom';
import Editor from '../Common/Editor';
import Flash from '../Common/Flash';
import { useDispatch, useSelector } from '../../store';
import { resetEditor } from '../../store/editor';
import { createBook } from '../../store/books';
import { setMessage } from '../../store/ui';
import PayWall from '../Common/PayWall';
import styles from './New.scss';

interface Props extends RouteComponentProps {}

// useInitFocus initializes the focus on HTML elements depending on the current
// state of the editor.
function useInitFocus({ bookLabel, textareaRef, setTriggerFocus }) {
  useEffect(() => {
    if (!bookLabel) {
      setTriggerFocus();
    } else {
      const textareaEl = textareaRef.current;

      if (textareaEl) {
        focusTextarea(textareaEl);
      }
    }
  }, [setTriggerFocus, bookLabel, textareaRef]);
}

const New: React.SFC<Props> = ({ history }) => {
  const { editor } = useSelector(state => {
    return {
      editor: state.editor
    };
  });
  const dispatch = useDispatch();
  const [errMessage, setErrMessage] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const textareaRef = useRef(null);
  const [setTriggerFocus, triggerRef] = useFocus();

  useCleanupEditor();
  useInitFocus({
    bookLabel: editor.bookLabel,
    textareaRef,
    setTriggerFocus
  });

  return (
    <Fragment>
      <Helmet>
        <title>New</title>
      </Helmet>

      <PayWall>
        <div
          className={classnames(
            styles.container,
            'container mobile-nopadding page page-mobile-full'
          )}
        >
          <Flash kind="danger" when={Boolean(errMessage)}>
            Error: {errMessage}
          </Flash>

          <div className={styles.wrapper}>
            <div className={classnames(styles.overlay, {})} />
            <div className={styles.header}>
              <h2 className={styles.heading}>New note</h2>
            </div>

            <Editor
              isNew
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
      </PayWall>
    </Fragment>
  );
};

export default React.memo(withRouter(New));
