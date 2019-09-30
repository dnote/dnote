import { BookData } from '../operations/books';

// Option represents an option in a selection list
export interface Option {
  label: string;
  value: string;
}

// optionValueCreate is the value of the option for creating a new option
export const optionValueCreate = 'create-new-option';

// filterOptions returns a new array of options based on the given filter criteria
export function filterOptions(
  options: Option[],
  term: string,
  creatable: boolean
): Option[] {
  if (!term) {
    return options;
  }

  const ret = [];
  const searchReg = new RegExp(`${term}`, 'i');
  let hit = null;

  for (let i = 0; i < options.length; i++) {
    const option = options[i];

    if (option.label === term) {
      hit = option;
    } else if (searchReg.test(option.label) && option.value !== '') {
      ret.push(option);
    }
  }

  // if there is an exact match, display the option at the top
  // otherwise, display a creatable option at the bottom
  if (hit) {
    ret.unshift(hit);
  } else if (creatable) {
    // creatable option has a value of an empty string
    ret.push({ label: term, value: '' });
  }

  return ret;
}

// booksToOptions returns an array of options for select ui, given an array of books
export function booksToOptions(books: BookData[]): Option[] {
  const ret = [];

  for (let i = 0; i < books.length; ++i) {
    const book = books[i];

    ret.push({
      label: book.label,
      value: book.uuid
    });
  }

  return ret;
}
