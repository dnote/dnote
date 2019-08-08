import React from 'react';
import classnames from 'classnames';

import styles from './SearchBar.scss';

interface Props {}

const SearchBar: React.SFC<Props> = () => {
  function handleSubmit(e) {
    e.preventDefault();
  }

  return (
    <form className={styles.wrapper} onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Search notes"
        className={classnames(styles.input, 'text-input text-input-small')}
      />
      <button type="submit" className={styles.button}>
        Search
      </button>
    </form>
  );
};

export default SearchBar;
