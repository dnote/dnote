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

import React, { useEffect, useState, Fragment } from 'react';
import classnames from 'classnames';

import { getBooks } from '../../store/books';
import { useDispatch } from '../../store';
import PayWall from '../Common/PayWall';
import Content from './Content';
import Flash from '../Common/Flash';
import HeadData from './HeadData';
import styles from './Books.scss';

const Books: React.SFC = () => {
  const [successMessage, setSuccessMessage] = useState('');

  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(getBooks());
  }, [dispatch]);

  return (
    <Fragment>
      <HeadData />

      <PayWall>
        <h1 className="sr-only">Books</h1>

        <div
          className={classnames('container mobile-nopadding', styles.wrapper)}
        >
          <Flash
            kind="success"
            when={Boolean(successMessage)}
            onDismiss={() => {
              setSuccessMessage('');
            }}
          >
            {successMessage}
          </Flash>

          <div className="row">
            <div className="col-12">
              <Content setSuccessMessage={setSuccessMessage} />;
            </div>
          </div>
        </div>
      </PayWall>
    </Fragment>
  );
};

export default Books;
