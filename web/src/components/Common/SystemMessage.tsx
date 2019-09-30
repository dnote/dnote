/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import React from 'react';
import classnames from 'classnames';
import { matchPath } from 'react-router';
import { withRouter, RouteComponentProps } from 'react-router-dom';
import { Location } from 'history';

import Flash from './Flash';
import { unsetMessage } from '../../store/ui';
import { useSelector, useDispatch } from '../../store';
import { MessageState } from '../../store/ui';
import styles from './SystemMessage.scss';

interface Props extends RouteComponentProps {}

function matchMessagePath(location: Location, message: MessageState): string {
  const paths = Object.keys(message);
  for (let i = 0; i < paths.length; ++i) {
    const path = paths[i];

    const match = matchPath(location.pathname, {
      path,
      exact: true
    });

    if (match) {
      return path;
    }
  }

  return null;
}

const SystemMessage: React.SFC<Props> = ({ location }) => {
  const { message } = useSelector(state => {
    return {
      message: state.ui.message
    };
  });
  const dispatch = useDispatch();

  const matchedPath = matchMessagePath(location, message);
  if (matchedPath === null) {
    return null;
  }

  const messageData = message[matchedPath];

  return (
    <div className={classnames('container mobile-nopadding', styles.wrapper)}>
      <Flash
        kind={messageData.kind}
        onDismiss={() => {
          dispatch(unsetMessage(matchedPath));
        }}
        noMargin
      >
        {messageData.content}
      </Flash>
    </div>
  );
};

export default withRouter(SystemMessage);
