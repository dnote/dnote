import {
  START_FETCHING,
  RECEIVE,
  RECEIVE_ERROR,
  BooksState,
  BooksActionType
} from './types';

const initialState = {
  items: [],
  isFetching: false,
  isFetched: false,
  error: null
};

export default function(
  state = initialState,
  action: BooksActionType
): BooksState {
  switch (action.type) {
    case START_FETCHING:
      return {
        ...state,
        isFetching: true,
        isFetched: false
      };
    case RECEIVE: {
      const { books } = action.data;

      // get uuids of deleted books and that of a currently selected book
      return {
        ...state,
        isFetching: false,
        isFetched: true,
        items: [...state.items, ...books]
      };
    }
    case RECEIVE_ERROR:
      return {
        ...state,
        isFetching: false,
        isFetched: true,
        error: action.data.error
      };
    default:
      return state;
  }
}
