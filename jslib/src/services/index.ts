/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import { HttpClientConfig } from '../helpers/http';
import initUsersService from './users';
import initBooksService from './books';
import initNotesService from './notes';
import initPaymentService from './payment';
import initDigestsService from './digests';
import initRepetitionRulesService from './repetitionRules';

// init initializes service helpers with the given http configuration
// and returns an object of all services.
export default function initServices(c: HttpClientConfig) {
  const usersService = initUsersService(c);
  const booksService = initBooksService(c);
  const notesService = initNotesService(c);
  const paymentService = initPaymentService(c);
  const digestsService = initDigestsService(c);
  const repetitionRulesService = initRepetitionRulesService(c);

  return {
    users: usersService,
    books: booksService,
    notes: notesService,
    payment: paymentService,
    digests: digestsService,
    repetitionRules: repetitionRulesService
  };
}
