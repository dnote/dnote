import { UPDATE, RESET, UpdateAction, ResetAction } from './types';

export function updateSettings(settings): UpdateAction {
  return {
    type: UPDATE,
    data: { settings }
  };
}

export function resetSettings(): ResetAction {
  return {
    type: RESET
  };
}
