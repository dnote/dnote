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

import React, { useState, useReducer, useRef, useEffect } from 'react';
import classnames from 'classnames';
import { Link } from 'react-router-dom';

import { getRepetitionsPath } from 'web/libs/paths';
import { Option, booksToOptions } from 'jslib/helpers/select';
import { BookDomain } from 'jslib/operations/types';
import { CreateParams } from 'jslib/services/repetitionRules';
import Modal, { Header, Body } from '../../Common/Modal';
import { useSelector } from '../../../store';
import { daysToMs } from '../../../helpers/time';
import Button from '../../Common/Button';
import MultiSelect from '../../Common/MultiSelect';
import styles from './Form.scss';
import modalStyles from '../../Common/Modal/Modal.scss';

export interface FormState {
  title: string;
  enabled: boolean;
  hour: number;
  minute: number;
  frequency: number;
  noteCount: number;
  bookDomain: BookDomain;
  books: Option[];
}

// serializeFormState serializes the given form state into a payload
export function serializeFormState(s: FormState): CreateParams {
  let bookUUIDs = [];
  if (s.bookDomain === BookDomain.All) {
    bookUUIDs = [];
  } else {
    bookUUIDs = s.books.map(b => {
      return b.value;
    });
  }

  return {
    title: s.title,
    hour: s.hour,
    minute: s.minute,
    frequency: s.frequency,
    book_domain: s.bookDomain,
    book_uuids: bookUUIDs,
    note_count: s.noteCount,
    enabled: s.enabled
  };
}

interface Props {
  onSubmit: (formState) => void;
  setErrMsg: (string) => void;
  cancelPath?: string;
  initialState?: FormState;
  isEditing?: boolean;
}

enum Action {
  setTitle,
  setFrequency,
  setHour,
  setMinutes,
  setNoteCount,
  setBookDomain,
  setBooks,
  toggleEnabled
}

function formReducer(state, action): FormState {
  switch (action.type) {
    case Action.setTitle:
      return {
        ...state,
        title: action.data
      };
    case Action.setFrequency:
      return {
        ...state,
        frequency: action.data
      };
    case Action.setHour:
      return {
        ...state,
        hour: action.data
      };
    case Action.setMinutes:
      return {
        ...state,
        minute: action.data
      };
    case Action.setNoteCount:
      return {
        ...state,
        noteCount: action.data
      };
    case Action.setBooks:
      return {
        ...state,
        books: action.data
      };
    case Action.setBookDomain:
      return {
        ...state,
        bookDomain: action.data
      };
    case Action.toggleEnabled:
      return {
        ...state,
        enabled: !state.enabled
      };
    default:
      return state;
  }
}

const formInitialState: FormState = {
  title: '',
  enabled: true,
  hour: 8,
  minute: 0,
  frequency: daysToMs(7),
  noteCount: 20,
  bookDomain: BookDomain.All,
  books: []
};

function validateForm(state: FormState): Error | null {
  if (state.title === '') {
    return new Error('Title is required.');
  }
  if (state.bookDomain !== BookDomain.All && state.books.length === 0) {
    return new Error('Please select books.');
  }
  if (state.noteCount <= 0) {
    return new Error('Please specify note count greater than 0.');
  }

  return null;
}

const Form: React.FunctionComponent<Props> = ({
  onSubmit,
  setErrMsg,
  cancelPath = getRepetitionsPath(),
  initialState = formInitialState,
  isEditing = false
}) => {
  const [inProgress, setInProgress] = useState(false);
  const bookSelectorInputRef = useRef(null);
  const [formState, formDispatch] = useReducer(formReducer, initialState);
  const { books } = useSelector(state => {
    return {
      books: state.books.data
    };
  });
  const bookOptions = booksToOptions(books);
  const booksSelectTextId = 'book-select-text-input';

  let bookSelectorPlaceholder;
  if (formState.bookDomain === BookDomain.All) {
    bookSelectorPlaceholder = 'All books';
  } else if (formState.bookDomain === BookDomain.Including) {
    bookSelectorPlaceholder = 'Select books to include';
  } else if (formState.bookDomain === BookDomain.Excluding) {
    bookSelectorPlaceholder = 'Select books to exclude';
  }

  let bookSelectorCurrentOptions;
  if (formState.bookDomain === BookDomain.All) {
    bookSelectorCurrentOptions = [];
  } else {
    bookSelectorCurrentOptions = formState.books;
  }

  useEffect(() => {
    if (isEditing) {
      return;
    }

    if (formState.bookDomain === BookDomain.All) {
      if (bookSelectorInputRef.current) {
        bookSelectorInputRef.current.blur();
      }
    } else {
      if (bookSelectorInputRef.current) {
        bookSelectorInputRef.current.focus();
      }
    }
  }, [formState.bookDomain, isEditing]);

  return (
    <form
      onSubmit={e => {
        e.preventDefault();

        const err = validateForm(formState);
        if (err !== null) {
          setErrMsg(err.message);
          return;
        }

        onSubmit(formState);
      }}
      className={styles.form}
    >
      <div className={styles['input-row']}>
        <label className="input-label" htmlFor="title">
          Name
        </label>

        <input
          autoFocus
          type="text"
          id="title"
          className="text-input text-input-small text-input-stretch"
          placeholder="Weekly vocabulary reminder"
          value={formState.title}
          onChange={e => {
            const data = e.target.value;

            formDispatch({
              type: Action.setTitle,
              data
            });
          }}
        />
      </div>

      <div className={styles['input-row']}>
        <label className="input-label" htmlFor={booksSelectTextId}>
          Eligible books
        </label>

        <div className={styles['book-domain-wrapper']}>
          <div className={styles['book-domain-option']}>
            <input
              type="radio"
              id="book-domain-all"
              name="book-domain"
              value="all"
              checked={formState.bookDomain === BookDomain.All}
              onChange={e => {
                const data = e.target.value;

                formDispatch({
                  type: Action.setBookDomain,
                  data
                });
              }}
            />
            <label
              className={styles['book-domain-label']}
              htmlFor="book-domain-all"
            >
              All
            </label>
          </div>

          <div className={styles['book-domain-option']}>
            <input
              type="radio"
              id="book-domain-including"
              name="book-domain"
              value="including"
              checked={formState.bookDomain === BookDomain.Including}
              onChange={e => {
                const data = e.target.value;

                formDispatch({
                  type: Action.setBookDomain,
                  data
                });
              }}
            />
            <label
              className={styles['book-domain-label']}
              htmlFor="book-domain-including"
            >
              Including
            </label>
          </div>

          <div className={styles['book-domain-option']}>
            <input
              type="radio"
              id="book-domain-excluding"
              name="book-domain"
              value="excluding"
              checked={formState.bookDomain === BookDomain.Excluding}
              onChange={e => {
                const data = e.target.value;

                formDispatch({
                  type: Action.setBookDomain,
                  data
                });
              }}
            />
            <label
              className={styles['book-domain-label']}
              htmlFor="book-domain-excluding"
            >
              Excluding
            </label>
          </div>
        </div>

        <MultiSelect
          disabled={formState.bookDomain === BookDomain.All}
          textInputId={booksSelectTextId}
          options={bookOptions}
          currentOptions={bookSelectorCurrentOptions}
          setCurrentOptions={data => {
            formDispatch({ type: Action.setBooks, data });
          }}
          placeholder={bookSelectorPlaceholder}
          wrapperClassName={styles['book-selector']}
          inputInnerRef={bookSelectorInputRef}
        />
      </div>

      <div
        className={classnames(styles['input-row'], styles['schedule-wrapper'])}
      >
        <div className={styles['schedule-content']}>
          <div className={classnames(styles['schedule-input-wrapper'])}>
            <label className="input-label" htmlFor="frequency">
              How often?
            </label>

            <select
              id="frequency"
              className="form-select"
              value={formState.frequency}
              onChange={e => {
                const { value } = e.target;

                formDispatch({
                  type: Action.setFrequency,
                  data: Number.parseInt(value)
                });
              }}
            >
              <option value={daysToMs(1)}>Every day</option>
              <option value={daysToMs(2)}>Every 2 days</option>
              <option value={daysToMs(3)}>Every 3 days</option>
              <option value={daysToMs(4)}>Every 4 days</option>
              <option value={daysToMs(5)}>Every 5 days</option>
              <option value={daysToMs(6)}>Every 6 days</option>
              <option value={daysToMs(7)}>Every week</option>
              <option value={daysToMs(14)}>Every 2 weeks</option>
              <option value={daysToMs(21)}>Every 3 weeks</option>
              <option value={daysToMs(28)}>Every 4 weeks</option>
            </select>
          </div>

          <div className={styles['schedule-input-wrapper']}>
            <label className="input-label" htmlFor="hour">
              Hour
            </label>

            <select
              id="hour"
              className={classnames('form-select', styles['time-select'])}
              value={formState.hour}
              onChange={e => {
                const { value } = e.target;

                formDispatch({
                  type: Action.setHour,
                  data: Number.parseInt(value, 10)
                });
              }}
            >
              {[...Array(24)].map((_, i) => {
                return (
                  <option key={i} value={i}>
                    {i}
                  </option>
                );
              })}
            </select>
          </div>

          <div className={styles['schedule-input-wrapper']}>
            <label className="input-label" htmlFor="minutes">
              Minutes
            </label>

            <select
              id="minutes"
              className={classnames('form-select', styles['time-select'])}
              value={formState.minute}
              onChange={e => {
                const { value } = e.target;

                formDispatch({
                  type: Action.setMinutes,
                  data: Number.parseInt(value, 10)
                });
              }}
            >
              {[...Array(60)].map((_, i) => {
                return (
                  <option key={i} value={i}>
                    {i}
                  </option>
                );
              })}
            </select>
          </div>
        </div>

        <div className={styles.help}>
          When to deliver a digest in the UTC (Coordinated Universal Time).
        </div>
      </div>

      <div className={styles['input-row']}>
        <label className="input-label" htmlFor="num-notes">
          Number of notes
        </label>

        <input
          type="number"
          min="1"
          id="num-notes"
          className="text-input text-input-small"
          placeholder="10"
          value={formState.noteCount}
          onChange={e => {
            const { value } = e.target;

            let data;
            if (value === '') {
              data = '';
            } else {
              data = Number.parseInt(value);
            }

            formDispatch({
              type: Action.setNoteCount,
              data
            });
          }}
        />

        <div className={styles.help}>
          Maximum number of notes to include in each repetition
        </div>
      </div>

      <div className={styles['input-row']}>
        <label className="input-label" htmlFor="enabled">
          Enabled?
        </label>

        <div>
          <input
            type="checkbox"
            id="enabled"
            checked={formState.enabled}
            onChange={e => {
              const data = e.target.value;

              formDispatch({
                type: Action.toggleEnabled,
                data
              });
            }}
          />
        </div>
      </div>

      <div className={modalStyles.actions}>
        <Button type="submit" kind="first" size="normal" isBusy={inProgress}>
          Create
        </Button>

        <Link
          to={cancelPath}
          onClick={e => {
            const ok = window.confirm('Are you sure?');
            if (!ok) {
              e.preventDefault();
              return;
            }
          }}
          className="button button-second button-normal"
        >
          Cancel
        </Link>
      </div>
    </form>
  );
};

export default Form;
