export interface EditorState {
  noteUUID: string | null;
  bookUUID: string | null;
  bookLabel: string | null;
  content: string;
  dirty: boolean;
}

export const MARK_DIRTY = 'editor/MARK_DIRTY';
export const STAGE_NOTE = 'editor/STAGE_NOTE';
export const FLUSH_CONTENT = 'editor/FLUSH_CONTENT';
export const UPDATE_BOOK = 'editor/UPDATE_BOOK';
export const RESET = 'editor/RESET';

export interface MarkDirtyAction {
  type: typeof MARK_DIRTY;
}

export interface StageNoteAction {
  type: typeof STAGE_NOTE;
  data: {
    noteUUID: string;
    bookUUID: string;
    bookLabel: string;
    content: string;
  };
}

export interface FlushContentAction {
  type: typeof FLUSH_CONTENT;
  data: {
    content: string;
  };
}

export interface UpdateBookAction {
  type: typeof UPDATE_BOOK;
  data: {
    uuid: string;
    label: string;
  };
}

export interface ResetAction {
  type: typeof RESET;
}

export type EditorActionType =
  | MarkDirtyAction
  | StageNoteAction
  | FlushContentAction
  | UpdateBookAction
  | ResetAction;
