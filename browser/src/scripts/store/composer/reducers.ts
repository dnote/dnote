import {
  UPDATE_CONTENT,
  UPDATE_BOOK,
  RESET,
  RESET_BOOK,
  ComposerActionType,
  ComposerState
} from './types';

const initialState: ComposerState = {
  content: '',
  bookUUID: '',
  bookLabel: ''
};

export default function(
  state = initialState,
  action: ComposerActionType
): ComposerState {
  switch (action.type) {
    case UPDATE_CONTENT: {
      return {
        ...state,
        content: action.data.content
      };
    }
    case UPDATE_BOOK: {
      return {
        ...state,
        bookUUID: action.data.uuid,
        bookLabel: action.data.label
      };
    }
    case RESET_BOOK: {
      return {
        ...state,
        bookUUID: '',
        bookLabel: ''
      };
    }
    case RESET: {
      return initialState;
    }
    default:
      return state;
  }
}
