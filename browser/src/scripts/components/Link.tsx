import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { useDispatch } from '../store/hooks';
import { navigate } from '../store/location/actions';

interface Props {
  to: string;
  className: string;
  tabIndex?: number;
  onClick?: () => void;
}

const Link: React.FunctionComponent<Props> = ({
  to,
  children,
  onClick,
  ...restProps
}) => {
  const dispatch = useDispatch();

  return (
    <a
      href={`${to}`}
      onClick={e => {
        e.preventDefault();

        dispatch(navigate(to));

        if (onClick) {
          onClick();
        }
      }}
      {...restProps}
    >
      {children}
    </a>
  );
};

export default Link;
