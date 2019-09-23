export interface SettingsState {
  sessionKey: string;
  sessionKeyExpiry: number;
}

export const UPDATE = 'settings/UPDATE';
export const RESET = 'settings/RESET';

export interface UpdateAction {
  type: typeof UPDATE;
  data: {
    settings: any;
  };
}

export interface ResetAction {
  type: typeof RESET;
}

export type SettingsActionType = UpdateAction | ResetAction;
