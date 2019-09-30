import React from 'react';
import classnames from 'classnames';

import styles from './AdvancedPanel.scss';

interface Props {
  words: string;
  setWords: (string) => void;
  disabled: boolean;
}

const WordsSearch: React.SFC<Props> = ({ words, setWords, disabled }) => {
  return (
    <section className={styles.section}>
      <label htmlFor="has-words" className={styles.label}>
        Has words
        <input
          type="text"
          id="has-words"
          className={classnames(
            'text-input text-input-small text-input-stretch',
            styles.input
          )}
          value={words}
          disabled={disabled}
          onChange={e => {
            const val = e.target.value;
            setWords(val);
          }}
        />
      </label>
    </section>
  );
};

export default WordsSearch;
