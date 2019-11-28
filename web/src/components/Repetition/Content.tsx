import React, { Fragment, useState, useEffect } from 'react';

import classnames from 'classnames';
import {
  getNewRepetitionPath,
  getSettingsPath,
  getSubscriptionPath,
  SettingSections,
  repetitionsPathDef
} from 'web/libs/paths';
import { Link } from 'react-router-dom';
import { useDispatch, useSelector } from '../../store';
import { getRepetitionRules } from '../../store/repetitionRules';
import RepetitionList from './RepetitionList';
import DeleteRepetitionRuleModal from './DeleteRepetitionRuleModal';
import Flash from '../Common/Flash';
import { setMessage } from '../../store/ui';
import styles from './Repetition.scss';

const Content: React.FunctionComponent = () => {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(getRepetitionRules());
  }, [dispatch]);

  const { repetitionRules, user } = useSelector(state => {
    return {
      repetitionRules: state.repetitionRules,
      user: state.auth.user.data
    };
  });

  const [ruleUUIDToDelete, setRuleUUIDToDelete] = useState('');

  return (
    <Fragment>
      <div className="container mobile-fw">
        <div className={classnames('page-header', styles.header)}>
          <h1 className="page-heading">Repetition</h1>

          {!user.pro ? (
            <button
              disabled
              type="button"
              className="button button-first button-normal"
            >
              New
            </button>
          ) : (
            <Link
              id="T-new-rule-btn"
              className="button button-first button-normal"
              to={getNewRepetitionPath()}
            >
              New
            </Link>
          )}
        </div>
      </div>

      <div className="container mobile-nopadding">
        <Flash when={!user.pro} kind="warning">
          Repetitions are not enabled on your plan.{' '}
          <Link to={getSubscriptionPath()}>Upgrade here.</Link>
        </Flash>

        <Flash when={user.pro && !user.emailVerified} kind="warning">
          Please verify your email address in order to receive digests.{' '}
          <Link to={getSettingsPath(SettingSections.account)}>
            Go to settings.
          </Link>
        </Flash>
        <RepetitionList
          isFetching={repetitionRules.isFetching}
          isFetched={repetitionRules.isFetched}
          items={repetitionRules.data}
          setRuleUUIDToDelete={setRuleUUIDToDelete}
          pro={user.pro}
        />
      </div>

      <DeleteRepetitionRuleModal
        repetitionRuleUUID={ruleUUIDToDelete}
        isOpen={ruleUUIDToDelete !== ''}
        onDismiss={() => {
          setRuleUUIDToDelete('');
        }}
        setSuccessMessage={message => {
          dispatch(
            setMessage({
              message,
              kind: 'info',
              path: repetitionsPathDef
            })
          );
        }}
      />
    </Fragment>
  );
};

export default Content;
