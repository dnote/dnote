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

import React, { Fragment, useRef, useState, useEffect } from 'react';

import { usePrevious } from '../../../libs/hooks';
import Editor from './Editor';
import Preview from './Preview';

// getEditorScrollHandler returns an event handler to roughly sync the vertical scroll
// position of editor and preview
function getEditorScrollHandler(previewEl) {
  let baseScrollTop = 0;
  let delta = 0;

  let ticking = false;

  // previewScale determines the scale of the scroll delta of the preview compared to
  // that of the editor. Huristically, it helps to keep scrolls in sync because
  // a rendered HTML typically takes more vertical spaces than the markdown.
  const previewScale = 1.5;

  function getUpdate() {
    // eslint-disable-next-line no-param-reassign
    previewEl.scrollTop += delta * previewScale;

    ticking = false;
    baseScrollTop = 0;
    delta = 0;
  }
  function requestTick() {
    if (ticking || delta === 0) {
      return;
    }

    window.requestAnimationFrame(getUpdate);
    ticking = true;
  }

  return e => {
    if (baseScrollTop === 0) {
      baseScrollTop = e.target.scrollTop;
    } else {
      delta = e.target.scrollTop - baseScrollTop;
    }

    requestTick();
  };
}

function useDraft(editorData) {
  const [draft, setDraft] = useState('');

  const prevNoteUUID = usePrevious(editorData.noteUUID);
  useEffect(() => {
    if (editorData.noteUUID !== prevNoteUUID) {
      setDraft(editorData.content);
    }
  }, [editorData.noteUUID, editorData.content, prevNoteUUID, setDraft]);

  return [draft, setDraft];
}

function Body({ editorData, textareaRef, noteData, demo, isReady }) {
  const [draft, setDraft] = useDraft(editorData);
  const previewRef = useRef(null);
  const onEditorScroll = getEditorScrollHandler(previewRef.current);

  return (
    <Fragment>
      <Editor
        draft={draft}
        setDraft={setDraft}
        editorData={editorData}
        textareaRef={textareaRef}
        isFetching={noteData.isFetching}
        disabled={!isReady}
        demo={demo}
        onPaneScroll={onEditorScroll}
      />
      <Preview
        noteError={noteData.error}
        content={draft}
        previewRef={previewRef}
      />
    </Fragment>
  );
}

export default Body;
