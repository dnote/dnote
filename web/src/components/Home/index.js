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

import React, { useEffect, useRef } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import Helmet from 'react-helmet';
import classnames from 'classnames';

import Workspace from './Workspace';

import { resetNote } from '../../actions/note';
import { resetEditor, stageNote } from '../../actions/editor';
import { getCipherKey } from '../../crypto';
import { getNote } from '../../actions/note';
import { usePrevious } from '../../libs/hooks';

import style from './Home.module.scss';

function Home({
  demo,
  editorData,
  match,
  doGetNote,
  doStageNote,
  doResetNote,
  doResetEditor
}) {
  const textareaRef = useRef(null);

  const { noteUUID } = match.params;
  const prevNoteUUID = usePrevious(noteUUID);
  useEffect(() => {
    if (!noteUUID) {
      return;
    }
    if (prevNoteUUID === noteUUID) {
      return;
    }

    const cipherKeyBuf = getCipherKey(demo);
    doGetNote(cipherKeyBuf, noteUUID, demo)
      .then(n => {
        // stage the note to the editor if editor is not dirty
        if (!editorData.dirty) {
          doStageNote({
            noteUUID: n.uuid,
            bookUUID: n.book.uuid,
            content: n.content
          });
        }

        const textareaEl = textareaRef.current;
        if (textareaEl) {
          textareaEl.focus();
        }
      })
      .catch(err => {
        console.log('error', err);
      });
  });

  useEffect(() => {
    return () => {
      doResetNote();
      doResetEditor();
    };
  }, [doResetNote, doResetEditor]);

  return (
    <div id="T-home-page" className={classnames(style.wrapper)}>
      <Helmet>
        <title>Notes</title>
      </Helmet>

      <Workspace textareaRef={textareaRef} demo={demo} />
    </div>
  );
}

function mapStateToProps(state) {
  return {
    note: state.note,
    layout: state.ui.layout,
    showNewNoteModal: state.ui.modal.newNote,
    editorData: state.editor
  };
}

const mapDispatchToProps = {
  doResetNote: resetNote,
  doResetEditor: resetEditor,
  doGetNote: getNote,
  doStageNote: stageNote
};

Home.defaultProps = {
  demo: false
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(Home)
);
