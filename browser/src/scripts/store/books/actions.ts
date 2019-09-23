import services from '../../utils/services';

import {
  START_FETCHING,
  RECEIVE,
  RECEIVE_ERROR,
  StartFetchingAction,
  ReceiveAction,
  ReceiveErrorAction
} from './types';

function startFetchingBooks(): StartFetchingAction {
  return {
    type: START_FETCHING
  };
}

function receiveBooks(books): ReceiveAction {
  return {
    type: RECEIVE,
    data: {
      books
    }
  };
}

function receiveBooksError(error: string): ReceiveErrorAction {
  return {
    type: RECEIVE_ERROR,
    data: {
      error
    }
  };
}

export function fetchBooks() {
  return (dispatch, getState) => {
    dispatch(startFetchingBooks());

    const { settings } = getState();

    services.books
      .fetch(
        {},
        {
          headers: {
            Authorization: `Bearer ${settings.sessionKey}`
          }
        }
      )
      .then(books => {
        dispatch(receiveBooks(books));
      })
      .catch(err => {
        console.log('error fetching books', err);
        dispatch(receiveBooksError(err));
      });
  };
}
