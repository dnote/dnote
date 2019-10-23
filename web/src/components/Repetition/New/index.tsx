import React, { useState, useEffect } from 'react';
import Helmet from 'react-helmet';
import { Link, withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';

import { getRepetitionsPath, repetitionsPathDef } from 'web/libs/paths';
import { BookDomain } from 'jslib/operations/types';
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
    </div>
  );
};

export default withRouter(NewRepetition);
