import React, { useState } from 'react';
import { connect } from 'react-redux';
import classnames from 'classnames';

import { KEYCODE_ENTER } from 'jslib/helpers/keyboard';
import { flushContent } from '../../../store/editor';
import { AppState } from '../../../store';
import styles from './Textarea.scss';
import editorStyles from './Editor.scss';

interface Props {
  content: string;
  onChange: (string) => void;
  doFlushContent: (string) => void;
  onSubmit: () => void;
  textareaRef: React.MutableRefObject<any>;
  inputTimerRef: React.MutableRefObject<any>;
  disabled?: boolean;
}

const Textarea: React.SFC<Props> = ({
  content,
  onChange,
  doFlushContent,
  onSubmit,
  textareaRef,
  inputTimerRef,
  disabled
}) => {
  const [contentFocused, setContentFocused] = useState(false);

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

            doFlushContent(value);
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

function mapStateToProps(state: AppState) {
  return {
    editor: state.editor
  };
}

const mapDispatchToProps = {
  doFlushContent: flushContent
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Textarea);
