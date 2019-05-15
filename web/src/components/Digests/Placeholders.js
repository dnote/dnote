import React from 'react';

import DigestHolder from './DigestHolder';

const placeholders = new Array(12);
for (let i = 0; i < placeholders.length; ++i) {
  placeholders[i] = <DigestHolder key={i} />;
}

function LoadingList() {
  return placeholders;
}

export default LoadingList;
