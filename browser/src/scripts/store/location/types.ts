export interface LocationState {
  path: string;
  state: any;
}

export const NAVIGATE = 'location/NAVIGATE';

export interface NavigateAction {
  type: typeof NAVIGATE;
  data: {
    path: string;
    state: string;
  };
}

export type LocationActionType = NavigateAction;
