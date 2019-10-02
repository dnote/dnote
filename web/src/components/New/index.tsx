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
import { RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';
import Helmet from 'react-helmet';
import { withRouter } from 'react-router-dom';

import { focusTextarea } from 'web/libs/dom';
import { getEditorSessionkey } from 'web/libs/editor';
import operations from 'web/libs/operations';
import { getNotePath, notePathDef } from 'web/libs/paths';
import { useFocus } from 'web/libs/hooks/dom';
import Editor from '../Common/Editor';
import Flash from '../Common/Flash';
import { useDispatch, useSelector } from '../../store';
import { resetEditor, createSession } from '../../store/editor';
import { createBook } from '../../store/books';
import { setMessage } from '../../store/ui';
import PayWall from '../Common/PayWall';
import Content from './Content';
import styles from './New.scss';

interface Props extends RouteComponentProps {}

// useInitFocus initializes the focus on HTML elements depending on the current
// state of the editor.
function useInitFocus({ bookLabel, content, textareaRef, setTriggerFocus }) {
  useEffect(() => {
    if (!bookLabel && !content) {
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
  const sessionKey = getEditorSessionkey(null);
  const { editor } = useSelector(state => {
    return {
      editor: state.editor
    };
  });

  const session = editor.sessions[sessionKey];

  const dispatch = useDispatch();
  useEffect(() => {
    // if there is no editorSesssion session, create one
    if (session === undefined) {
      dispatch(
        createSession({
          noteUUID: null,
          bookUUID: null,
          bookLabel: null,
          content: ''
        })
      );
    }
  }, [dispatch, session]);

  return (
    <Fragment>
      <Helmet>
        <title>New</title>
      </Helmet>

      {session !== undefined && (
        <Content editor={session} persisted={editor.persisted} />
      )}
    </Fragment>
  );
};

export default React.memo(withRouter(New));
