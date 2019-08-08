export interface MessageData {
  content: string;
  kind: string;
}

export interface MessageState {
  [path: string]: MessageData;
}

export interface UIState {
  message: MessageState;
}

export const SET_MESSAGE = 'ui/SET_MESSAGE';
export const UNSET_MESSAGE = 'ui/UNSET_MESSAGE';

export interface SetMessageAction {
  type: typeof SET_MESSAGE;
  data: {
    message: string;
    kind: string;
    path: string;
  };
}

export interface UnsetMessageAction {
  type: typeof UNSET_MESSAGE;
  data: {
    path: string;
  };
}

export type UIActionType = SetMessageAction | UnsetMessageAction;
