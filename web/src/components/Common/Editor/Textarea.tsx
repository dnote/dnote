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

import React, { useState } from 'react';
import classnames from 'classnames';

import { KEYCODE_ENTER } from 'jslib/helpers/keyboard';
import { flushContent } from '../../../store/editor';
import { AppState, useDispatch } from '../../../store';
import styles from './Textarea.scss';
import editorStyles from './Editor.scss';

interface Props {
  sessionKey: string;
  content: string;
  onChange: (string) => void;
  onSubmit: () => void;
  textareaRef: React.MutableRefObject<any>;
  inputTimerRef: React.MutableRefObject<any>;
  disabled?: boolean;
}

const Textarea: React.SFC<Props> = ({
  sessionKey,
  content,
  onChange,
  onSubmit,
  textareaRef,
  inputTimerRef,
  disabled
}) => {
  const [contentFocused, setContentFocused] = useState(false);
  const dispatch = useDispatch();

  return (
    <div className={classnames(styles.wrapper, editorStyles.content)}>
      <textarea
        ref={textareaRef}
        value={content}
        onChange={e => {
          const { value } = e.target;
          onChange(value);

          // flush the draft to the data store when user stops typing
          if (inputTimerRef.current) {
            window.clearTimeout(inputTimerRef.current);
          }
          // eslint-disable-next-line no-param-reassign
          inputTimerRef.current = window.setTimeout(() => {
            // eslint-disable-next-line no-param-reassign
            inputTimerRef.current = null;

            dispatch(flushContent(sessionKey, value));
          }, 1000);
        }}
        onFocus={() => {
          setContentFocused(true);
        }}
        onKeyDown={e => {
          if (e.shiftKey && e.keyCode === KEYCODE_ENTER) {
            e.preventDefault();

            onSubmit();
          }
        }}
        onBlur={() => setContentFocused(false)}
        className={classnames(styles.textarea, 'text-input')}
        placeholder="What did you learn?"
        disabled={disabled}
      />

      <span
        className={classnames(styles.tip, { [styles.shown]: contentFocused })}
      >
        Shift + Enter to save
      </span>
    </div>
  );
};

export default Textarea;
