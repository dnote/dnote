import React, { useState, useCallback } from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import * as queriesLib from 'jslib/helpers/queries';
import { getSearchDest } from 'web/libs/search';
import { useFilters } from '../../../../store';
import Button from '../../../Common/Button';
import PopoverContent from '../../../Common/Popover/PopoverContent';
import BookSearch from './BookSearch';
import WordsSearch from './WordsSearch';
import styles from './AdvancedPanel.scss';

interface Props extends RouteComponentProps {
  onDismiss: () => void;
  disabled: boolean;
}

// quoteFilters surrounds a filter term with a pair of double quotation marks, effectively
// making it a text term.
function quoteFilters(s: string): string {
  let ret = '';

  const terms = s.split(' ');
  for (let i = 0; i < terms.length; ++i) {
    const term = terms[i];

    if (i > 0) {
      ret += ' ';
    }

    const parts = term.split(':');

    if (parts.length === 2 && queriesLib.keywords.indexOf(parts[0]) > -1) {
      ret += `"${term}"`;
    } else {
      ret += `${term}`;
    }
  }

  return ret;
}

const quotedRegex = /"(.*)"/;

// unquoteFilters removes surrounding double quotation marks for a valid filter term
function unquoteFilters(s: string): string {
  let ret = '';

  const terms = s.split(' ');
  for (let i = 0; i < terms.length; ++i) {
    const term = terms[i];

    if (i > 0) {
      ret += ' ';
    }

    const matchGroup = term.match(quotedRegex);
    if (matchGroup !== null) {
      const match = matchGroup[1];
      const parts = match.split(':');

      if (parts.length === 2 && queriesLib.keywords.indexOf(parts[0]) > -1) {
        ret += match;
      }
    } else {
      ret += term;
    }
  }

  return ret;
}

function encodeBookStr(s: string): string[] {
  const ret = [];

  const parts = s.split(',');
  for (let i = 0; i < parts.length; i++) {
    const p = parts[i];
    const candidate = p.trim();

    if (candidate !== '') {
      ret.push(candidate);
    }
  }

  return ret;
}

const AdvancedPanel: React.SFC<Props> = ({
  onDismiss,
  history,
  location,
  disabled
}) => {
  const filters = useFilters();
  const { queries } = filters;

  const [words, setWords] = useState(unquoteFilters(queries.q));
  const [books, setBooks] = useState(queries.book.join(', '));

  const handleSubmit = useCallback(() => {
    const q: queriesLib.Queries = {
      q: quoteFilters(words),
      book: encodeBookStr(books)
    };

    const dest = getSearchDest(location, q);
    history.push(dest);

    onDismiss();
  }, [history, onDismiss, location, words, books]);

  return (
    <PopoverContent
      wrapperEl={document}
      contentId="advanced-search-panel"
      onDismiss={onDismiss}
      contentClassName={styles.wrapper}
      direction="bottom"
      closeOnEscapeKeydown
      closeOnOutsideClick
    >
      <form className={styles.form} onSubmit={handleSubmit}>
        <WordsSearch words={words} setWords={setWords} disabled={disabled} />

        <BookSearch value={books} setValue={setBooks} disabled={disabled} />

        <Button
          type="submit"
          kind="first"
          size="normal"
          stretch
          className={styles.submit}
          disabled={disabled}
        >
          Search
        </Button>
      </form>
    </PopoverContent>
  );
};

export default withRouter(AdvancedPanel);
