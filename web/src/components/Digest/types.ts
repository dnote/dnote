// Sort is a set of possible values for sort query parameters
export enum Sort {
  Newest = '',
  Oldest = 'created-asc'
}

export interface SearchParams {
  sort: Sort;
  books: string[];
}
