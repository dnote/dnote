/* eslint-disable react/no-danger */

import React from 'react';
import classnames from 'classnames';

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
          className={classnames('markdown-body', styles.content)}
          dangerouslySetInnerHTML={{
            __html: formatContent(note.content)
          }}
        />
      )}
    </section>
  );
};

export default Content;
