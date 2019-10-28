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

import React, { useState, useEffect } from 'react';
import classnames from 'classnames';
import { Link, withRouter, RouteComponentProps } from 'react-router-dom';
import Helmet from 'react-helmet';

import services from 'web/libs/services';
import Logo from '../../Icons/Logo';
import Flash from '../../Common/Flash';
import { parseSearchString } from 'jslib/helpers/url';
import { getEmailPreference } from '../../../store/auth';
import { getLoginPath } from 'web/libs/paths';
import { useSelector, useDispatch } from '../../../store';
import Content from './Content';
import styles from './EmailPreferenceRepetition.scss';

interface Match {
  repetitionUUID: string;
}
interface Props extends RouteComponentProps<Match> {}

const EmailPreferenceRepetition: React.SFC<Props> = ({ location, match }) => {
  const [data, setData] = useState(null);
  const [isFetching, setIsFetching] = useState(false);
  const [successMsg, setSuccessMsg] = useState('');
  const [failureMsg, setFailureMsg] = useState('');

  const { token } = parseSearchString(location.search);
  const { repetitionUUID } = match.params;

  useEffect(() => {
    if (data !== null) {
      return;
    }

    setIsFetching(true);

    services.repetitionRules
      .fetch(repetitionUUID, { token })
      .then(repetition => {
        setData(repetition);
        setIsFetching(false);
      })
      .catch(err => {
        if (err.response.status === 401) {
          setFailureMsg('Your email token has expired or is not valid.');
        } else {
          setFailureMsg(err.message);
        }

        setIsFetching(false);
      });
  }, [data, setData, setFailureMsg, setIsFetching]);

  const isFetched = data !== null;

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>Toggle repetition</title>
      </Helmet>

      <Link to="/">
        <Logo fill="#252833" width={60} height={60} />
      </Link>
      <h1 className={styles.heading}>Toggle repetition</h1>

      <div className="container">
        <div className={styles.body}>
          <Flash
            when={failureMsg !== ''}
            kind="danger"
            wrapperClassName={styles.flash}
          >
            {failureMsg}{' '}
            <span>
              Please <Link to={getLoginPath()}>login</Link> and try again.
            </span>
          </Flash>

          <Flash
            when={successMsg !== ''}
            kind="success"
            wrapperClassName={classnames(styles.flash, 'T-success')}
          >
            {successMsg}
          </Flash>

          {isFetching && <div>Loading</div>}

          {isFetched && (
            <Content
              token={token}
              data={data}
              setSuccessMsg={setSuccessMsg}
              setFailureMsg={setFailureMsg}
            />
          )}
        </div>
        <div className={styles.footer}>
          <Link to="/">Back to Dnote home</Link>
        </div>
      </div>
    </div>
  );
};

export default withRouter(EmailPreferenceRepetition);
