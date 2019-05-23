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
import classnames from 'classnames';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import DotsIcon from '../../Icons/Dots';
import Menu from '../../Common/Menu';
import { getHomePath } from '../../../libs/paths';

import { removeNote } from '../../../actions/notes';
import { resetNote } from '../../../actions/note';
import { resetEditor } from '../../../actions/editor';

import styles from './NoteActions.module.scss';
import workspaceStyle from './Workspace.module.scss';
import * as notesOperation from '../../../operations/notes';

function handleRemove({
  event,
  noteUUID,
  note,
  doRemoveNote,
  doResetNote,
  doResetEditor,
  history
}) {
  event.preventDefault();

  const ok = window.confirm('Are you sure to remove this note?');

  if (!ok) {
    return;
  }

  notesOperation
    .remove(noteUUID)
    .then(() => {
      console.log('remove successful');

      const date = new Date(note.created_at);
      const year = date.getUTCFullYear();
      const month = date.getUTCMonth() + 1;

      // TODO: focus the next note

      doRemoveNote({ year, month, noteUUID });
      doResetNote();
      doResetEditor();

      console.log('pushing');
      history.push(getHomePath());
    })
    .catch(err => {
      console.log('err', err);
    });
}

function NoteActions({
  noteUUID,
  noteData,
  doRemoveNote,
  doResetNote,
  doResetEditor,
  history,
  demo,
  disabled
}) {
  const [isOpen, setIsOpen] = useState(false);
  const optRefs = [useRef(null)];

  const options = [
    {
      name: 'Remove',
      value: (
        <form
          onSubmit={event => {
            setIsOpen(false);
            handleRemove({
              event,
              noteUUID,
              note: noteData.item,
              doRemoveNote,
              doResetNote,
              doResetEditor,
              history
            });
          }}
        >
          <button
            id="T-remove-note-btn"
            role="menuitem"
            type="submit"
            className={classnames(
              'button-no-ui button-stretch',
              styles.action,
              {
                disabled: demo
              }
            )}
            disabled={demo}
            ref={optRefs[0]}
          >
            Remove
          </button>
        </form>
      )
    }
  ];

  return (
    <Menu
      triggerId="note-actions"
      options={options}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      optRefs={optRefs}
      triggerContent={<DotsIcon width={16} height={16} />}
      triggerClassName={workspaceStyle.action}
      contentClassName={styles.content}
      wrapperClassName={styles.wrapper}
      alignment="right"
      direction="bottom"
      disabled={disabled}
    />
  );
}

function mapStateToProps(state) {
  return {
    noteData: state.note
  };
}

const mapDispatchToProps = {
  doRemoveNote: removeNote,
  doResetNote: resetNote,
  doResetEditor: resetEditor
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(NoteActions)
);
