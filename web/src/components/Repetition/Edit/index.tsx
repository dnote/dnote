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
import Helmet from 'react-helmet';
import { Link, withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';

import { getRepetitionsPath, repetitionsPathDef } from 'web/libs/paths';
import { BookDomain, RepetitionRuleData } from 'jslib/operations/types';
import services from 'web/libs/services';
import { createRepetitionRule } from '../../../store/repetitionRules';
import { useDispatch } from '../../../store';
import Flash from '../../Common/Flash';
import { setMessage } from '../../../store/ui';
import Content from './Content';
import repetitionStyles from '../Repetition.scss';

interface Match {
  repetitionUUID: string;
}

interface Props extends RouteComponentProps<Match> {}

const EditRepetition: React.FunctionComponent<Props> = ({ history, match }) => {
  const dispatch = useDispatch();
  const [errMsg, setErrMsg] = useState('');
  const [data, setData] = useState<RepetitionRuleData | null>(null);

  useEffect(() => {
    const { repetitionUUID } = match.params;
    services.repetitionRules
      .fetch(repetitionUUID)
      .then(rule => {
        setData(rule);
      })
      .catch(err => {
        setErrMsg(err.message);
      });
  }, [dispatch, match]);

  return (
    <div id="page-edit-repetition" className="page page-mobile-full">
      <Helmet>
        <title>Edit Repetition</title>
      </Helmet>

      <div className="container mobile-fw">
        <div className={classnames('page-header', repetitionStyles.header)}>
          <h1 className="page-heading">Edit Repetition</h1>

          <Link to={getRepetitionsPath()}>Back</Link>
        </div>

        <Flash
          kind="danger"
          when={errMsg !== ''}
          onDismiss={() => {
            setErrMsg('');
          }}
        >
          Error: {errMsg}
        </Flash>

        {data === null ? (
          <div>loading</div>
        ) : (
          <Content setErrMsg={setErrMsg} data={data} />
        )}
      </div>
    </div>
  );
};

export default EditRepetition;
