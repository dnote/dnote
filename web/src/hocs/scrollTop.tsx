import React, { useEffect } from 'react';
import { withRouter, RouteComponentProps } from 'react-router-dom';

interface Props extends RouteComponentProps {}

export default function(WrappedComponent: React.ComponentType) {
  const ScrollToTop: React.SFC<Props> = ({ location }) => {
    const { pathname } = location;
    useEffect(() => {
      console.log('scrolling to top');
      window.scrollTo(0, 0);
    }, [pathname]);

    return <WrappedComponent />;
  };

  return withRouter(ScrollToTop);
}
