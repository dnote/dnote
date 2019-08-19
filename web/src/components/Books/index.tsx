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

import React, { useEffect, useState } from 'react';

import { getBooks } from '../../store/books';
import { useDispatch, useSelector } from '../../store';
import SubscriberWall from '../Common/SubscriberWall';
import Content from './Content';
import Flash from '../Common/Flash';
import HeadData from './HeadData';

interface ContentWrapperProps {
  setSuccessMessage: (string) => void;
}

const ContentWrapper: React.SFC<ContentWrapperProps> = ({
  setSuccessMessage
}) => {
  const { user } = useSelector(state => {
    return {
      user: state.auth.user
    };
  });

  if (user.data.pro) {
    return <Content setSuccessMessage={setSuccessMessage} />;
  }

  return <SubscriberWall />;
};

const Books: React.SFC = () => {
  const [successMessage, setSuccessMessage] = useState('');

  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(getBooks());
  }, [dispatch]);

  return (
    <div>
      <HeadData />

      {/* <h1 className="sr-only">Books</h1> */}

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

        <div className="row">
          <div className="col-12">
            <ContentWrapper setSuccessMessage={setSuccessMessage} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Books;
