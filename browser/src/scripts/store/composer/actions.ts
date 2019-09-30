import {
  UPDATE_CONTENT,
  UPDATE_BOOK,
  RESET,
  RESET_BOOK,
  UpdateContentAction,
  UpdateBookAction,
  ResetBookAction,
  ResetAction
} from './types';

export function updateContent(content: string): UpdateContentAction {
  return {
    type: UPDATE_CONTENT,
    data: { content }
  };
}

export interface UpdateBookActionParam {
  uuid: string;
  label: string;
}

export function updateBook({
  uuid,
  label
}: UpdateBookActionParam): UpdateBookAction {
  return {
    type: UPDATE_BOOK,
    data: {
      uuid,
      label
    }
  };
}

export function resetBook(): ResetBookAction {
  return {
    type: RESET_BOOK
  };
}

export function resetComposer(): ResetAction {
  return {
    type: RESET
  };
}
