import { NAVIGATE, LocationState, LocationActionType } from './types';

const initialState: LocationState = {
  path: '/',
  state: {}
};

export default function(
  state = initialState,
  action: LocationActionType
): LocationState {
  switch (action.type) {
    case NAVIGATE:
      return {
        ...state,
        path: action.data.path,
        state: action.data.state || {}
      };
    default:
      return state;
  }
}
