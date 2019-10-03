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

import { getBooks } from '../../store/books';
import { useDispatch } from '../../store';
import PayWall from '../Common/PayWall';
import Content from './Content';
import Flash from '../Common/Flash';
import HeadData from './HeadData';

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
        <div className="page page-mobile-full">
          <div className="container mobile-nopadding">
            <Flash
              kind="success"
              when={Boolean(successMessage)}
              onDismiss={() => {
                setSuccessMessage('');
              }}
            >
              {successMessage}
            </Flash>
          </div>
          <Content setSuccessMessage={setSuccessMessage} />;
        </div>
      </PayWall>
    </Fragment>
  );
};

export default Books;
