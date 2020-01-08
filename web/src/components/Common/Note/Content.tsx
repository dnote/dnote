/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

/* eslint-disable react/no-danger */

import React from 'react';

import { NoteData } from 'jslib/operations/types';
import { excerpt } from 'web/libs/string';
import { tokenize, TokenKind } from 'web/libs/fts/lexer';
import { parseMarkdown } from '../../../helpers/markdown';
import styles from './Note.scss';

function formatFTSSelection(content: string): string {
  if (content.indexOf('<dnotehl>') === -1) {
    return content;
  }

  const tokens = tokenize(content);

  let output = '';
  let buf = [];

  for (let i = 0; i < tokens.length; i++) {
    const t = tokens[i];

    if (t.kind === TokenKind.hlBegin || t.kind === TokenKind.eol) {
      output += buf.join('');

      buf = [];
    } else if (t.kind === TokenKind.hlEnd) {
      output += `<span class="${styles.match}">
        ${buf.join('')}
      </span>`;

      buf = [];
    } else {
      buf.push(t.value);
    }
  }

  return output;
}

function formatContent(content: string): string {
  const formatted = formatFTSSelection(content);
  return parseMarkdown(formatted);
}

interface Props {
  collapsed?: boolean;
  note: NoteData;
}

const Content: React.SFC<Props> = ({ note, collapsed }) => {
  return (
    <section className={styles['content-wrapper']}>
      {collapsed ? (
        <div className={styles['collapsed-content']}>
          {excerpt(note.content, 100)}
        </div>
      ) : (
        <div
          className="markdown-body"
          dangerouslySetInnerHTML={{
            __html: formatContent(note.content)
          }}
        />
      )}
    </section>
  );
};

export default Content;
