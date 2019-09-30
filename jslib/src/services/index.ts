import { HttpClientConfig } from '../helpers/http';
import initUsersService from './users';
import initBooksService from './books';
import initNotesService from './notes';
import initDigestsService from './digests';
import initPaymentService from './payment';

// init initializes service helpers with the given http configuration
// and returns an object of all services.
export default function initServices(c: HttpClientConfig) {
  const usersService = initUsersService(c);
  const booksService = initBooksService(c);
  const notesService = initNotesService(c);
  const digestsService = initDigestsService(c);
  const paymentService = initPaymentService(c);

  return {
    users: usersService,
    books: booksService,
    notes: notesService,
    digests: digestsService,
    payment: paymentService
  };
}
