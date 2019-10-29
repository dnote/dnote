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
import { Link } from 'react-router-dom';
import classnames from 'classnames';

import { getNewRepetitionPath } from 'web/libs/paths';
import { getRepetitionRules } from '../../store/repetitionRules';
import PayWall from '../Common/PayWall';
import { useDispatch, useSelector } from '../../store';
import RepetitionList from './RepetitionList';
import DeleteRepetitionRuleModal from './DeleteRepetitionRuleModal';
import Flash from '../Common/Flash';
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

  const [ruleUUIDToDelete, setRuleUUIDToDelete] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

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
              id="T-new-rule-btn"
              className="button button-first button-normal"
              to={getNewRepetitionPath()}
            >
              New
            </Link>
          </div>
        </div>

        <div className="container mobile-nopadding">
          <Flash
            when={successMessage !== ''}
            kind="success"
            onDismiss={() => {
              setSuccessMessage('');
            }}
          >
            {successMessage}
          </Flash>

          <RepetitionList
            isFetching={repetitionRules.isFetching}
            isFetched={repetitionRules.isFetched}
            items={repetitionRules.data}
            setRuleUUIDToDelete={setRuleUUIDToDelete}
          />
        </div>

        <DeleteRepetitionRuleModal
          repetitionRuleUUID={ruleUUIDToDelete}
          isOpen={ruleUUIDToDelete !== ''}
          onDismiss={() => {
            setRuleUUIDToDelete('');
          }}
          setSuccessMessage={setSuccessMessage}
        />
      </PayWall>
    </div>
  );
};

export default Repetition;
