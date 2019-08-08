import React, { useEffect, useState } from 'react';
import { Prompt, RouteComponentProps } from 'react-router-dom';
import Helmet from 'react-helmet';
import { withRouter } from 'react-router-dom';

import Flash from '../Common/Flash';
import { useDispatch, useSelector } from '../../store';
import { stageNote } from '../../store/editor';
import * as notesOperation from '../../operations/notes';
import Content from './Content';

interface Match {
  noteUUID: string;
}

interface Props extends RouteComponentProps<Match> {}

const Edit: React.SFC<Props> = ({ match }) => {
  const { editor } = useSelector(state => {
    return {
      editor: state.editor
    };
  });
  const dispatch = useDispatch();

  const [errMessage, setErrMessage] = useState('');
  const [isReady, setIsReady] = useState(false);

  const { noteUUID } = match.params;

  useEffect(() => {
    notesOperation
      .fetchOne(noteUUID)
      .then(note => {
        dispatch(
          stageNote({
            noteUUID: note.uuid,
            bookUUID: note.book.uuid,
            bookLabel: note.book.label,
            content: note.content
          })
        );

        setIsReady(true);
      })
      .catch((err: Error) => {
        setErrMessage(err.message);
      });
  }, [dispatch, noteUUID]);

  return (
    <div className="container mobile-nopadding">
      <Helmet>
        <title>Edit Note</title>
      </Helmet>

      <Flash kind="danger" when={Boolean(errMessage)}>
        Error: {errMessage}
      </Flash>

      <div className="row">
        <div className="col-12">
          {isReady ? (
            <Content noteUUID={noteUUID} setErrMessage={setErrMessage} />
          ) : null}
        </div>
      </div>

      <Prompt
        message="You have unsaved changes. Continue?"
        when={editor.dirty}
      />
    </div>
  );
};

export default React.memo(withRouter(Edit));
