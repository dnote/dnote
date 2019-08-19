import React from 'react';
import Helmet from 'react-helmet';

import { Filters } from '../../libs/filters';

interface Props {
  filters: Filters;
}

function getTitle(filters: Filters): string {
  if (filters.queries.book.length === 1) {
    return `Notes in ${filters.queries.book}`;
  }

  return 'Notes';
}

const HeaderData: React.SFC<Props> = ({ filters }) => {
  const title = getTitle(filters);

  return (
    <Helmet>
      <title>{title}</title>
    </Helmet>
  );
};

export default HeaderData;
