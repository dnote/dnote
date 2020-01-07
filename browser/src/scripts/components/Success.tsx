/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import React, { useEffect, useState, Fragment } from 'react';

import {
  KEYCODE_ENTER,
  KEYCODE_ESC,
  KEYCODE_LOWERCASE_B
} from 'jslib/helpers/keyboard';
import Flash from './Flash';
import ext from '../utils/ext';
import BookIcon from './BookIcon';
import { navigate } from '../store/location/actions';
import { useSelector, useDispatch } from '../store/hooks';

const Success: React.FunctionComponent = () => {
  const [errorMsg, setErrorMsg] = useState('');

  const dispatch = useDispatch();
  const { location, settings } = useSelector(state => ({
    location: state.location,
    settings: state.settings
  }));

  const { bookName, noteUUID } = location.state;

  useEffect(() => {
    const handleKeydown = e => {
      e.preventDefault();

      if (e.keyCode === KEYCODE_ENTER) {
        dispatch(navigate('/'));
      } else if (e.keyCode === KEYCODE_ESC) {
        window.close();
      } else if (e.keyCode === KEYCODE_LOWERCASE_B) {
        const url = `${settings.webUrl}/notes/${noteUUID}`;

        ext.tabs
          .create({ url })
          .then(() => {
            window.close();
          })
          .catch(err => {
            setErrorMsg(err.message);
          });
      }
    };

    window.addEventListener('keydown', handleKeydown);

    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  }, [dispatch, noteUUID, settings.webUrl]);

  return (
    <Fragment>
      <Flash kind="error" when={errorMsg !== ''} message={errorMsg} />

      <div className="success-page">
        <div>
          <BookIcon width={20} height={20} className="book-icon" />

          <h1 className="heading">
            Saved to
            {bookName}
          </h1>
        </div>

        <ul className="key-list">
          <li className="key-item">
            <kbd className="key">Enter</kbd>{' '}
            <div className="key-desc">Go back</div>
          </li>
          <li className="key-item">
            <kbd className="key">b</kbd>{' '}
            <div className="key-desc">Open in browser</div>
          </li>
          <li className="key-item">
            <kbd className="key">ESC</kbd>
            <div className="key-desc">Close</div>
          </li>
        </ul>
      </div>
    </Fragment>
  );
};

export default Success;
