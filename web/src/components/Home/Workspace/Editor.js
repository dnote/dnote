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
import classnames from 'classnames';
import { connect } from 'react-redux';

import Placeholder from './EditorPlaceholder';
import { updateContent, commitNote, markDirty } from '../../../actions/editor';

import workspaceStyle from './Workspace.module.scss';
import styles from './Editor.module.scss';

// useCommitNote periodically persists the current editor state to the service
function useCommitNote({ editorData, doCommitNote, demo }) {
  const timerRef = useRef(null);

  useEffect(() => {
    if (demo) {
      return;
    }
    if (!editorData.dirty) {
      return;
    }

    if (timerRef.current) {
      window.clearTimeout(timerRef.current);
    }

    timerRef.current = window.setTimeout(() => {
      doCommitNote();
    }, 1000);
  });
}

function Body({
  draft,
  setDraft,
  textareaRef,
  disabled,
  doUpdateContent,
  doCommitNote,
  doMarkDirty,
  editorData,
  demo,
  onPaneScroll
}) {
  useCommitNote({ editorData, doCommitNote, demo });
  const inputTimerRef = useRef(null);

  return (
    <textarea
      id="editor"
      className={classnames(styles.textarea, workspaceStyle['pane-content'])}
      value={draft}
      placeholder="Start typing..."
      onScroll={onPaneScroll}
      onChange={e => {
        const { value } = e.target;
        setDraft(value);

        if (!demo && !editorData.dirty) {
          doMarkDirty();
        }

        // flush the draft to the data store when user stops typing
        if (inputTimerRef.current) {
          window.clearTimeout(inputTimerRef.current);
        }
        inputTimerRef.current = window.setTimeout(() => {
          inputTimerRef.current = null;

          doUpdateContent(value);
        }, 1000);
      }}
      ref={textareaRef}
      disabled={disabled}
    />
  );
}

function Editor({
  draft,
  setDraft,
  textareaRef,
  disabled,
  doUpdateContent,
  doCommitNote,
  doMarkDirty,
  editorData,
  demo,
  isFetching,
  onPaneScroll
}) {
  return (
    <div className={classnames(workspaceStyle.pane, styles.wrapper)}>
      {isFetching ? (
        <Placeholder />
      ) : (
        <Body
          draft={draft}
          setDraft={setDraft}
          textareaRef={textareaRef}
          disabled={disabled}
          doUpdateContent={doUpdateContent}
          doCommitNote={doCommitNote}
          doMarkDirty={doMarkDirty}
          editorData={editorData}
          demo={demo}
          isFetching={isFetching}
          onPaneScroll={onPaneScroll}
        />
      )}
    </div>
  );
}

const mapDispatchToProps = {
  doUpdateContent: updateContent,
  doCommitNote: commitNote,
  doMarkDirty: markDirty
};

export default connect(
  null,
  mapDispatchToProps
)(Editor);
