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
