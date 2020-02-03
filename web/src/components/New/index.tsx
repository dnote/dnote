/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import React, { Fragment, useEffect } from 'react';
import Helmet from 'react-helmet';
import { RouteComponentProps, withRouter } from 'react-router-dom';
import { getEditorSessionkey } from 'web/libs/editor';
import { useDispatch, useSelector } from '../../store';
import { createSession } from '../../store/editor';
import PayWall from '../Common/PayWall';
import Content from './Content';

interface Props extends RouteComponentProps {}

const New: React.FunctionComponent<Props> = () => {
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
        <PayWall>
          <Content editor={session} persisted={editor.persisted} />
        </PayWall>
      )}
    </Fragment>
  );
};

export default React.memo(withRouter(New));
