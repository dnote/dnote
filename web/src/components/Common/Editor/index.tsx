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
import { Link } from 'react-router-dom';
import { Location } from 'history';

import { focusTextarea } from 'web/libs/dom';
import { getHomePath } from 'web/libs/paths';
import BooksSelector from './BookSelector';
import { useDispatch, useSelector } from '../../../store';
import {
  flushContent,
  resetEditor,
  EditorSession
} from '../../../store/editor';
import Textarea from './Textarea';
import Preview from './Preview';
import Button from '../Button';
import styles from './Editor.scss';

function getContentCacheKey(editorSessionKey: string) {
  return `editor.${editorSessionKey}.content`;
}

function useEditorContent(
  editor: EditorSession,
  cacheKey: string
): [string, React.Dispatch<any>] {
  const cached = localStorage.getItem(cacheKey);
  return useState(cached || editor.content);
}

interface Props {
  editor: EditorSession;
  onSubmit: (param: { draftContent: string; draftBookUUID: string }) => void;
  isBusy: boolean;
  cancelPath?: Location<any>;
  isNew?: boolean;
  disabled?: boolean;
  textareaRef: React.MutableRefObject<any>;
  bookSelectorTriggerRef?: React.MutableRefObject<HTMLElement>;
}

enum Mode {
  write,
  preview
}

const Editor: React.SFC<Props> = ({
  editor,
  onSubmit,
  isBusy,
  disabled,
  textareaRef,
  isNew,
  bookSelectorTriggerRef,
  cancelPath = getHomePath()
}) => {
  const { books } = useSelector(state => {
    return {
      books: state.books
    };
  });
  const dispatch = useDispatch();
  const [bookSelectorOpen, setBookSelectorOpen] = useState(false);

  const contentCacheKey = getContentCacheKey(editor.sessionKey);
  const [content, setContent] = useEditorContent(editor, contentCacheKey);
  const [mode, setMode] = useState(Mode.write);
  const inputTimerRef = useRef(null);

  const isWriteMode = mode === Mode.write;
  const isPreviewMode = mode === Mode.preview;

  function handleSubmit() {
    // immediately flush the content
    if (inputTimerRef.current) {
      window.clearTimeout(inputTimerRef.current);

      // eslint-disable-next-line no-param-reassign
      inputTimerRef.current = null;
      dispatch(flushContent(editor.sessionKey, content));
    }

    onSubmit({ draftContent: content, draftBookUUID: editor.bookUUID });
    localStorage.removeItem(contentCacheKey);
  }

  if (disabled) {
    return <div>loading</div>;
  }

  return (
    <form
      id="T-editor"
      className={styles.wrapper}
      onSubmit={e => {
        e.preventDefault();
        handleSubmit();
      }}
    >
      <div className={classnames(styles.row, styles['editor-header'])}>
        <div>
          <BooksSelector
            editor={editor}
            isReady={books.isFetched}
            isOpen={bookSelectorOpen}
            setIsOpen={setBookSelectorOpen}
            triggerRef={bookSelectorTriggerRef}
            onAfterChange={() => {
              if (textareaRef.current) {
                focusTextarea(textareaRef.current);
              }
            }}
          />
        </div>

        <nav className={styles.tabs}>
          <button
            type="button"
            role="tab"
            aria-selected={isWriteMode}
            className={classnames('button-no-ui', styles.tab, {
              [styles['tab-active']]: isWriteMode
            })}
            onClick={() => {
              setMode(Mode.write);
            }}
          >
            Write
          </button>

          <button
            type="button"
            role="tab"
            aria-selected={isPreviewMode}
            className={classnames('button-no-ui', styles.tab, {
              [styles['tab-active']]: isPreviewMode
            })}
            onClick={() => {
              setMode(Mode.preview);
            }}
          >
            Preview
          </button>
        </nav>
      </div>

      <div className={styles['content-wrapper']}>
        {mode === Mode.write ? (
          <Textarea
            sessionKey={editor.sessionKey}
            textareaRef={textareaRef}
            inputTimerRef={inputTimerRef}
            content={content}
            onChange={c => {
              localStorage.setItem(contentCacheKey, c);
              setContent(c);
            }}
            onSubmit={handleSubmit}
          />
        ) : (
          <Preview content={content} />
        )}
      </div>

      <div className={styles.actions}>
        <Button
          id="T-save-button"
          type="submit"
          kind="third"
          size="normal"
          disabled={isBusy}
        >
          {isNew ? 'Save' : 'Update'}
        </Button>

        <Link
          to={cancelPath}
          onClick={e => {
            const ok = window.confirm('Are you sure?');
            if (!ok) {
              e.preventDefault();
              return;
            }

            localStorage.removeItem(contentCacheKey);
            dispatch(resetEditor(editor.sessionKey));
          }}
          className="button button-second button-normal"
        >
          Cancel
        </Link>
      </div>
    </form>
  );
};

export default Editor;
