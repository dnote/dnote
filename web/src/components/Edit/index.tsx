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

import React, { useEffect, useState } from 'react';
import classnames from 'classnames';
import { Prompt, RouteComponentProps } from 'react-router-dom';
import Helmet from 'react-helmet';
import { withRouter } from 'react-router-dom';

import { getEditorSessionkey } from 'web/libs/editor';
import operations from 'web/libs/operations';
import Flash from '../Common/Flash';
import { useDispatch, useSelector } from '../../store';
import { createSession } from '../../store/editor';
import Content from './Content';
import styles from '../New/New.scss';

interface Match {
  noteUUID: string;
}

interface Props extends RouteComponentProps<Match> {}

const Edit: React.SFC<Props> = ({ match }) => {
  const { noteUUID } = match.params;

  const sessionKey = getEditorSessionkey(noteUUID);
  const { editor } = useSelector(state => {
    return {
      editor: state.editor
    };
  });
  const session = editor.sessions[sessionKey];

  const dispatch = useDispatch();
  const [errMessage, setErrMessage] = useState('');

  useEffect(() => {
    if (session === undefined) {
      operations.notes
        .fetchOne(noteUUID)
        .then(note => {
          dispatch(
            createSession({
              noteUUID: note.uuid,
              bookUUID: note.book.uuid,
              bookLabel: note.book.label,
              content: note.content
            })
          );
        })
        .catch((err: Error) => {
          setErrMessage(err.message);
        });
    }
  }, [dispatch, noteUUID, session]);

  return (
    <div
      id="T-edit-page"
      className={classnames(
        styles.container,
        'container mobile-nopadding page page-mobile-full'
      )}
    >
      <Helmet>
        <title>Edit Note</title>
      </Helmet>

      <Flash kind="danger" when={Boolean(errMessage)}>
        Error: {errMessage}
      </Flash>

      {session !== undefined && (
        <Content
          noteUUID={noteUUID}
          editor={session}
          persisted={editor.persisted}
          setErrMessage={setErrMessage}
        />
      )}
    </div>
  );
};

export default React.memo(withRouter(Edit));
