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

import React, { useCallback } from 'react';
import moment from 'moment';
import { withRouter } from 'react-router-dom';
import classnames from 'classnames';
import { connect } from 'react-redux';

import SafeLink from '../Common/Link/SafeLink';
import { closeNoteSidebar } from '../../actions/ui';
import { notePath } from '../../libs/paths';
import { nanosecToSec } from '../../helpers/time';
import { excerpt } from '../../libs/string';
import { getWindowWidth, noteSidebarThreshold } from '../../libs/ui';
import { parseSearchString } from '../../libs/url';
import styles from './NoteItem.module.scss';

// renderContent renders the first line of the note
function renderContent(content) {
  let linebreakIdx = content.indexOf('\n');

  if (linebreakIdx === -1) {
    linebreakIdx = content.indexOf('\r\n');
  }

  let firstline;
  if (linebreakIdx === -1) {
    firstline = content;
  } else {
    firstline = content.substr(0, linebreakIdx);
  }

  return excerpt(firstline, 70);
}

function renderError() {
  return <div>Could not display the note</div>;
}

function renderBody({
  noteContent,
  noteUUID,
  editorContent,
  editorNoteUUID,
  isSelected,
  errorMessage
}) {
  if (errorMessage) {
    return renderError(errorMessage);
  }

  let content;
  // If note is selected and is loaded to the editor, display the editor content in the preview
  if (isSelected && editorNoteUUID === noteUUID) {
    content = editorContent;
  } else {
    content = noteContent;
  }

  return renderContent(content);
}

function NoteItem({
  note,
  demo,
  location,
  isSelected,
  errorMessage,
  editorData,
  doCloseNoteSidebar
}) {
  const queryObj = parseSearchString(location.search);

  const handleNoteItemClick = useCallback(() => {
    const width = getWindowWidth();

    if (width < noteSidebarThreshold) {
      doCloseNoteSidebar();
    }
  }, [doCloseNoteSidebar]);

  return (
    <li
      className={classnames('T-note-item', styles.wrapper, {
        [styles.active]: isSelected,
        [styles.error]: errorMessage
      })}
    >
      <SafeLink
        className={styles.link}
        to={notePath(note.uuid, queryObj, { demo, isEditor: true })}
        draggable="false"
        onClick={handleNoteItemClick}
      >
        {!errorMessage && (
          <div className={styles.meta}>
            <div className={styles['book-label']}>{note.book.label}</div>

            <div className={styles.ts}>
              {moment.unix(nanosecToSec(note.added_on)).fromNow()}
            </div>
          </div>
        )}
        <div className={styles.content}>
          {renderBody({
            noteContent: note.content,
            noteUUID: note.uuid,
            editorContent: editorData.content,
            editorNoteUUID: editorData.noteUUID,
            isSelected,
            errorMessage
          })}
        </div>
      </SafeLink>
    </li>
  );
}

function mapStateToProps(state) {
  return {
    editorData: state.editor
  };
}

const mapDispatchToProps = {
  doCloseNoteSidebar: closeNoteSidebar
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(NoteItem)
);
