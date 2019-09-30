import { UPDATE, RESET, SettingsState, SettingsActionType } from './types';

const initialState: SettingsState = {
  sessionKey: '',
  sessionKeyExpiry: 0
};

export default function(
  state = initialState,
  action: SettingsActionType
): SettingsState {
  switch (action.type) {
    case UPDATE:
      return {
        ...state,
        ...action.data.settings
      };
    case RESET:
      return initialState;
    default:
      return state;
  }
}
