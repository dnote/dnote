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
import React, { useEffect, useState } from 'react';
import Helmet from 'react-helmet';
import { RouteComponentProps, withRouter } from 'react-router-dom';
import { getEditorSessionkey } from 'web/libs/editor';
import operations from 'web/libs/operations';
import { useDispatch, useSelector } from '../../store';
import { createSession } from '../../store/editor';
import Flash from '../Common/Flash';
import styles from '../New/New.scss';
import Content from './Content';

interface Match {
  noteUUID: string;
}

interface Props extends RouteComponentProps<Match> {}

const Edit: React.FunctionComponent<Props> = ({ match }) => {
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
