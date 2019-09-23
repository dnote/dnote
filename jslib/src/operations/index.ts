import { HttpClientConfig } from '../helpers/http';
import initBooksOperation from './books';
import initNotesOperation from './notes';

// init initializes operations with the given http configuration
// and returns an object of all services.
export default function initOperations(c: HttpClientConfig) {
  const booksOperation = initBooksOperation(c);
  const notesOperation = initNotesOperation(c);

  return {
    books: booksOperation,
    notes: notesOperation
  };
}
