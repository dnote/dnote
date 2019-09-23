import React, { Fragment, useEffect, useState } from 'react';

import {
  KEYCODE_ENTER,
  KEYCODE_ESC,
  KEYCODE_LOWERCASE_B
} from 'jslib/helpers/keyboard';
import Flash from './Flash';
import ext from '../utils/ext';
import config from '../utils/config';
import BookIcon from './BookIcon';
import { navigate } from '../store/location/actions';
import { useSelector, useDispatch } from '../store/hooks';

const Success: React.FunctionComponent = () => {
  const [errorMsg, setErrorMsg] = useState('');

  const dispatch = useDispatch();
  const { location } = useSelector(state => {
    return {
      location: state.location
    };
  });

  const { bookName, noteUUID } = location.state;

  const handleKeydown = e => {
    e.preventDefault();

    if (e.keyCode === KEYCODE_ENTER) {
      dispatch(navigate('/'));
    } else if (e.keyCode === KEYCODE_ESC) {
      window.close();
    } else if (e.keyCode === KEYCODE_LOWERCASE_B) {
      const url = `${config.webUrl}/notes/${noteUUID}`;

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

  useEffect(() => {
    window.addEventListener('keydown', handleKeydown);

    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  }, []);

  return (
    <Fragment>
      <Flash when={errorMsg !== ''} message={errorMsg} />

      <div className="success-page">
        <div>
          <BookIcon width={20} height={20} className="book-icon" />

          <h1 className="heading">Saved to {bookName}</h1>
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
            <kbd className="key">ESC</kbd> <div className="key-desc">Close</div>
          </li>
        </ul>
      </div>
    </Fragment>
  );
};

export default Success;
