export interface QueryState {
  q: string;
  book: string[];
}

export interface FiltersState {
  queries: QueryState;
  page: number;
}

export const UPDATE_PAGE = 'filters/UDPATE_PAGE';
export const UPDATE_QUERY = 'filters/UPDATE_QUERY';
export const RESET = 'filters/RESET';

type ValidQueries = 'q' | 'book';

export interface UpdateQueryAction {
  type: typeof UPDATE_QUERY;
  data: {
    key: ValidQueries;
    value: string;
  };
}

export interface UpdatePageAction {
  type: typeof UPDATE_PAGE;
  data: {
    value: number;
  };
}

export interface ResetAction {
  type: typeof RESET;
}

export type FilterActionType =
  | UpdateQueryAction
  | UpdatePageAction
  | ResetAction;
