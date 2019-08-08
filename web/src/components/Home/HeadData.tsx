import React from 'react';
import Helmet from 'react-helmet';

interface Props {}

const HeaderData: React.SFC<Props> = () => {
  return (
    <Helmet>
      <title>Home</title>
    </Helmet>
  );
};

export default HeaderData;
