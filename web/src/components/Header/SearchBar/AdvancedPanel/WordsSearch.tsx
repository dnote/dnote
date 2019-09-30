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
