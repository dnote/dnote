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

import React, { useRef, useCallback, useState, useEffect } from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';

import * as filtersLib from 'jslib/helpers/filters';
import * as queriesLib from 'jslib/helpers/queries';
import { getSearchDest } from 'web/libs/search';
import { usePrevious } from 'web/libs/hooks';
import { useFilters, useSelector } from '../../../store';
import SearchInput from '../../Common/SearchInput';
import AdvancedPanel from './AdvancedPanel';
import styles from './SearchBar.scss';

const searchDelay = 930;

interface Props extends RouteComponentProps {}

const SearchBar: React.FunctionComponent<Props> = ({ location, history }) => {
  const searchTimerRef = useRef(null);
  const filters = useFilters();

  const initialValue = queriesLib.stringify(filters.queries);
  const [value, setValue] = useState(initialValue);
  const [expanded, setExpanded] = useState(false);

  const handleSearch = useCallback(
    (queryText: string) => {
      const queries = queriesLib.parse(queryText);
      const dest = getSearchDest(location, queries);
      history.push(dest);
    },
    [history, location]
  );

  const prevFilters = usePrevious(filters);
  useEffect(() => {
    if (prevFilters && filtersLib.checkFilterEqual(filters, prevFilters)) {
      return () => null;
    }

    const newVal = queriesLib.stringify(filters.queries);
    setValue(newVal);

    return () => null;
  }, [prevFilters, filters]);

  const onDismiss = () => {
    setExpanded(false);
  };

  const { user } = useSelector(state => {
    return {
      user: state.auth.user.data
    };
  });

  return (
    <div className={styles.wrapper}>
      <form
        onSubmit={e => {
          e.preventDefault();

          handleSearch(value);
        }}
        className={styles.form}
      >
        <SearchInput
          placeholder="Search notes"
          wrapperClassName={styles['input-wrapper']}
          inputClassName={classnames(styles.input, ' text-input-small')}
          disabled={!user.pro}
          value={value}
          onChange={e => {
            const val = e.target.value;
            setValue(val);

            if (searchTimerRef.current) {
              window.clearTimeout(searchTimerRef.current);
            }
            searchTimerRef.current = window.setTimeout(() => {
              handleSearch(val);
            }, searchDelay);
          }}
          onReset={() => {
            if (searchTimerRef.current) {
              window.clearTimeout(searchTimerRef.current);
            }

            handleSearch('');
          }}
          expanded={expanded}
          setExpanded={setExpanded}
        />
        <button type="submit" className={styles.button}>
          Search
        </button>
      </form>

      {expanded && <AdvancedPanel onDismiss={onDismiss} disabled={!user.pro} />}
    </div>
  );
};

export default withRouter(SearchBar);
