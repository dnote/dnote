import React, { useState, useEffect } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import { getNewRepetitionPath } from 'web/libs/paths';
import { getRepetitionRules } from '../../store/repetitionRules';
import PayWall from '../Common/PayWall';
import { useDispatch, useSelector } from '../../store';
import RepetitionList from './RepetitionList';
import styles from './Repetition.scss';

const Repetition: React.FunctionComponent = () => {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(getRepetitionRules());
  }, [dispatch]);

  const { repetitionRules } = useSelector(state => {
    return {
      repetitionRules: state.repetitionRules
    };
  });

  return (
    <div className="page page-mobile-full">
      <Helmet>
        <title>Repetition</title>
      </Helmet>

      <PayWall>
        <div className="container mobile-fw">
          <div className={classnames('page-header', styles.header)}>
            <h1 className="page-heading">Repetition</h1>

            <Link
              className="button button-first button-normal"
              to={getNewRepetitionPath()}
            >
              New
            </Link>
          </div>
        </div>

        <div className="container">
          <RepetitionList
            isFetching={repetitionRules.isFetching}
            isFetched={repetitionRules.isFetched}
            items={repetitionRules.data}
          />
        </div>
      </PayWall>
    </div>
  );
};

export default Repetition;
