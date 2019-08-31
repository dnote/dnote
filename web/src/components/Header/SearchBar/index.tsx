import React, { useRef, useCallback, useState, useEffect } from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';

import * as filtersLib from '../../../libs/filters';
import * as queriesLib from '../../../libs/queries';
import { usePrevious } from '../../../libs/hooks';
import { useFilters, useSelector } from '../../../store';
import SearchInput from '../../Common/SearchInput';
import AdvancedPanel from './AdvancedPanel';
import styles from './SearchBar.scss';

const searchDelay = 930;

interface Props extends RouteComponentProps {}

const SearchBar: React.SFC<Props> = ({ location, history }) => {
  const searchTimerRef = useRef(null);
  const filters = useFilters();

  const initialValue = queriesLib.stringify(filters.queries);
  const [value, setValue] = useState(initialValue);
  const [expanded, setExpanded] = useState(false);

  const { user } = useSelector(state => {
    return {
      user: state.auth.user.data
    };
  });

  const handleSearch = useCallback(
    (queryText: string) => {
      if (!user.pro) {
        return;
      }

      const queries = queriesLib.parse(queryText);
      const dest = queriesLib.getSearchDest(location, queries);
      history.push(dest);
    },
    [history, location, user]
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

  const disabled = !user.pro;

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
          value={value}
          disabled={disabled}
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

      {expanded && <AdvancedPanel onDismiss={onDismiss} disabled={disabled} />}
    </div>
  );
};

export default withRouter(SearchBar);
