import React from 'react';

import FeatureItem from './FeatureItem';

import styles from './FeatureList.module.scss';

function FeatureList({ features }) {
  return (
    <ul className={styles['feature-list']}>
      {features.map(feature => {
        return <FeatureItem key={feature.id} label={feature.label} />;
      })}
    </ul>
  );
}

export default FeatureList;
