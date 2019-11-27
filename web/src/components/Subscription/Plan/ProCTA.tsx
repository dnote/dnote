import React from 'react';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import {
  getSubscriptionCheckoutPath,
  getSettingsPath,
  SettingSections
} from 'web/libs/paths';
import { UserData } from 'jslib/operations/types';

interface Props {
  user: UserData;
}

const ProCTA: React.FunctionComponent<Props> = ({ user }) => {
  if (user && user.pro) {
    return (
      <Link
        to={getSettingsPath(SettingSections.billing)}
        className="button button-large button-third-outline button-stretch"
      >
        Manage Your Plan
      </Link>
    );
  }

  return (
    <Link
      id="T-unlock-pro-btn"
      className={classnames('button button-large button-third button-stretch')}
      to={getSubscriptionCheckoutPath()}
    >
      Upgrade
    </Link>
  );
};

export default ProCTA;
