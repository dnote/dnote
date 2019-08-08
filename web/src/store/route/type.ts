export interface RouteState {
  prevLocation: {
    pathname: string;
    search: string;
    hash: string;
    state: object;
  };
}

export const SET_PREV_LOCATION = 'route/SET_PREV_LOCATION';

export interface SetPrevLocationAction {
  type: typeof SET_PREV_LOCATION;
  data: {
    pathname: string;
    search: string;
    hash: string;
    state?: object;
  };
}

export type RouteActionType = SetPrevLocationAction;
