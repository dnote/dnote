/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
import PayWall from '../../Common/PayWall';
import {
  getRepetitionRules,
  createRepetitionRule
} from '../../../store/repetitionRules';
import { useDispatch } from '../../../store';
import Form, { FormState, serializeFormState } from '../Form';
import Flash from '../../Common/Flash';
import { setMessage } from '../../../store/ui';
import repetitionStyles from '../Repetition.scss';

interface Props extends RouteComponentProps {}

const NewRepetition: React.FunctionComponent<Props> = ({ history }) => {
  const dispatch = useDispatch();
  const [errMsg, setErrMsg] = useState('');

  useEffect(() => {
    dispatch(getRepetitionRules());
  }, [dispatch]);

  async function handleSubmit(state: FormState) {
    const payload = serializeFormState(state);

    try {
      await dispatch(createRepetitionRule(payload));

      const dest = getRepetitionsPath();
      history.push(dest);

      dispatch(
        setMessage({
          message: 'Created a repetition rule',
          kind: 'info',
          path: repetitionsPathDef
        })
      );
    } catch (e) {
      console.log(e);
      setErrMsg(e.message);
    }
  }

  return (
    <div id="page-new-repetition" className="page page-mobile-full">
      <Helmet>
        <title>New Repetition</title>
      </Helmet>

      <PayWall>
        <div className="container mobile-fw">
          <div className={classnames('page-header', repetitionStyles.header)}>
            <h1 className="page-heading">New Repetition</h1>

            <Link to={getRepetitionsPath()}>Back</Link>
          </div>

          <Flash
            kind="danger"
            when={errMsg !== ''}
            onDismiss={() => {
              setErrMsg('');
            }}
          >
            Error creating a rule: {errMsg}
          </Flash>

          <Form onSubmit={handleSubmit} setErrMsg={setErrMsg} />
        </div>
      </PayWall>
    </div>
  );
};

export default withRouter(NewRepetition);
