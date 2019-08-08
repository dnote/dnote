import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import classnames from 'classnames';

import { flushContent } from '../../../store/editor';
import { KEYCODE_ENTER } from '../../../helpers/keyboard';
import { AppState } from '../../../store';
import styles from './Textarea.scss';
import editorStyles from './Editor.scss';

interface Props {
  content: string;
  onChange: (string) => void;
  doFlushContent: (string) => void;
  onSubmit: () => void;
  setTextareaEl: React.Dispatch<any>;
  inputTimerRef: React.MutableRefObject<any>;
  disabled?: boolean;
}

const Textarea: React.SFC<Props> = ({
  content,
  onChange,
  doFlushContent,
  onSubmit,
  setTextareaEl,
  inputTimerRef,
  disabled
}) => {
  const [contentFocused, setContentFocused] = useState(false);

  useEffect(() => {
    return () => {
      // eslint-disable-next-line no-param-reassign
      setTextareaEl(null);
    };
  }, [setTextareaEl]);

  return (
    <div className={classnames(styles.wrapper, editorStyles.content)}>
      <textarea
        ref={el => {
          setTextareaEl(el);
        }}
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
