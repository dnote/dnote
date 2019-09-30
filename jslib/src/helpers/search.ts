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

export enum TokenKind {
  colon = 'COLON',
  id = 'ID',
  eof = 'EOF'
}

interface Token {
  kind: TokenKind;
  value?: string;
}

export enum NodeKind {
  text = 'text',
  filter = 'filter'
}

interface TextNode {
  kind: NodeKind.text;
  value: string;
}

interface FilterNode {
  kind: NodeKind.filter;
  keyword: string;
  value: string;
}

type Node = TextNode | FilterNode;

function nodeToString(node: Node): string {
  if (node.kind === NodeKind.text) {
    return node.value;
  }
  if (node.kind === NodeKind.filter) {
    return `${node.keyword}:${node.value}`;
  }

  throw new Error('unknown node kind');
}

const whitespaceRegex = /^\s*$/;
const charRegex = /^((?![\s:]).)*$/;

function isSpace(c: string): boolean {
  return whitespaceRegex.test(c);
}

function isChar(c: string): boolean {
  return charRegex.test(c);
}

export function tokenize(s: string): Token[] {
  const ret: Token[] = [];
  let cursor = 0;

  function colon() {
    ret.push({ kind: TokenKind.colon });
    cursor++;
  }
  function id() {
    let text = '';
    while (cursor <= s.length - 1 && isChar(s[cursor])) {
      text += s[cursor];
      cursor++;
    }

    ret.push({ kind: TokenKind.id, value: text });
  }

  while (cursor <= s.length - 1) {
    const currentChar = s[cursor];

    if (isSpace(currentChar)) {
      cursor++;
    } else if (currentChar === ':') {
      colon();
    } else if (isChar(currentChar)) {
      id();
    } else {
      throw new Error(`invalid character ${currentChar}`);
    }
  }

  ret.push({ kind: TokenKind.eof });

  return ret;
}

class Parser {
  toks: Token[];
  cursor: number;
  currentToken: Token;

  constructor(s: string) {
    this.toks = tokenize(s);
    this.cursor = 0;
  }

  do(): Node[] {
    /*
     * expr: term (term)*
     *
     * term: filter | text
     *
     * filter: text COLON text
     *
     * text: ID | EOF
     */
    return this.expr();
  }

  eat(kind: TokenKind) {
    const currentToken = this.getCurrentToken();

    if (currentToken.kind !== kind) {
      throw new Error(
        `invalid syntax. Expected ${kind} got: ${currentToken.kind}`
      );
    }

    this.cursor++;
  }

  getCurrentToken() {
    return this.toks[this.cursor];
  }

  expr(): Node[] {
    const nodes: Node[] = [];

    while (this.getCurrentToken().kind !== TokenKind.eof) {
      const n = this.term();
      nodes.push(n);
    }

    return nodes;
  }

  term(): Node | null {
    // try to parse filter and backtrack if not match
    const n = this.filter();
    if (n !== null) {
      return n;
    }

    return this.str();
  }

  filter(): Node | null {
    // save the current cursor for backtracking
    const cursor = this.cursor;

    const keyword = this.text({ maybe: true });
    if (keyword === null) {
      this.cursor = cursor;
      return null;
    }
    if (this.getCurrentToken().kind === TokenKind.colon) {
      this.eat(TokenKind.colon);
    } else {
      this.cursor = cursor;
      return null;
    }
    const value = this.text({ maybe: true });
    if (value === null) {
      this.cursor = cursor;
      return null;
    }

    return {
      kind: NodeKind.filter,
      keyword,
      value
    };
  }

  str(): Node | null {
    const value = this.text();
    if (value === null) {
      return null;
    }

    return {
      kind: NodeKind.text,
      value
    };
  }

  // text parses and returns a text. If 'maybe' option is true, it returns null
  // if the current token cannot be parsed as a text (such scenario can happen when
  // backtracking)
  text(opts = { maybe: false }): string | null {
    const currentToken = this.getCurrentToken();

    if (currentToken.kind === TokenKind.eof) {
      return null;
    }

    if (opts.maybe) {
      if (currentToken.kind !== TokenKind.id) {
        return null;
      }
    }

    this.eat(TokenKind.id);

    return currentToken.value;
  }
}

interface Search {
  text: string;
  filters: {
    [key: string]: string | string[];
  };
}

export function parse(s: string, keywords: string[]): Search {
  const p = new Parser(s);

  const ret: Search = {
    text: '',
    filters: {}
  };

  function addText(t: string) {
    if (ret.text !== '') {
      ret.text += ' ';
    }
    ret.text += t;
  }

  function addFilter(key: string, val: string) {
    const currentVal = ret.filters[key];

    if (typeof currentVal === 'undefined') {
      ret.filters[key] = val;
    } else if (typeof currentVal === 'string') {
      ret.filters[key] = [currentVal, val];
    } else {
      ret.filters[key] = [...currentVal, val];
    }
  }

  const nodes = p.do();
  for (let i = 0; i < nodes.length; i++) {
    const n = nodes[i];

    if (n.kind === NodeKind.text) {
      addText(n.value);
    } else if (n.kind === NodeKind.filter) {
      // if keyword was not specified, treat as a text
      if (keywords.indexOf(n.keyword) > -1) {
        addFilter(n.keyword, n.value);
      } else {
        addText(nodeToString(n));
      }
    }
  }

  return ret;
}
