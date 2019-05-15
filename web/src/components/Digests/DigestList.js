import React from 'react';

import DigestItem from './DigestItem';
import Placeholders from './Placeholders';

import styles from './DigestList.module.scss';

function Digests({ items, demo }) {
  return items.map(item => {
    return <DigestItem digest={item} demo={demo} />;
  });
}

function DigestList({ items, isFetching, demo }) {
  return (
    <ul className={styles.list}>
      {isFetching ? <Placeholders /> : <Digests items={items} demo={demo} />}
    </ul>
  );
}

export default DigestList;
