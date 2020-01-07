/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import React from 'react';
import CreatableSelect from 'react-select/creatable';
import cloneDeep from 'lodash/cloneDeep';
import { useSelector, useDispatch } from '../store/hooks';
import { updateBook, resetBook } from '../store/composer/actions';

interface Props {
  selectorRef: React.Dispatch<any>;
  onAfterChange: () => void;
}

function useCurrentOptions(options) {
  const currentValue = useSelector(state => state.composer.bookUUID);

  for (let i = 0; i < options.length; i++) {
    const option = options[i];

    if (option.value === currentValue) {
      return option;
    }
  }

  return null;
}

function useOptions() {
  const { books, composer } = useSelector(state => ({
    books: state.books,
    composer: state.composer
  }));

  const opts = books.items.map(book => ({
    label: book.label,
    value: book.uuid
  }));

  if (composer.bookLabel !== '' && composer.bookUUID === '') {
    opts.push({
      label: composer.bookLabel,
      value: ''
    });
  }

  // clone the array so as not to mutate Redux state manually
  // e.g. react-select mutates options prop internally upon adding a new option
  return cloneDeep(opts);
}

const BookSelector: React.FunctionComponent<Props> = ({
  selectorRef,
  onAfterChange
}) => {
  const dispatch = useDispatch();
  const { books } = useSelector(state => ({
    books: state.books
  }));
  const options = useOptions();
  const currentOption = useCurrentOptions(options);

  let placeholder: string;
  if (books.isFetched) {
    placeholder = 'Choose a book';
  } else {
    placeholder = 'Loading books...';
  }

  return (
    <CreatableSelect
      ref={el => {
        selectorRef(el);
      }}
      multi={false}
      isClearable
      placeholder={placeholder}
      options={options}
      value={currentOption}
      onChange={(option, meta) => {
        if (meta.action === 'clear') {
          dispatch(resetBook());
        } else {
          let uuid: string;
          if (meta.action === 'create-option') {
            uuid = '';
          } else {
            uuid = option.value;
          }

          dispatch(updateBook({ uuid, label: option.label }));
        }

        onAfterChange();
      }}
      formatCreateLabel={label => `Add a new book ${label}`}
      isDisabled={!books.isFetched}
    />
  );
};

export default BookSelector;
