import React, { useState, useEffect, useRef } from 'react';
import classnames from 'classnames';

import { KEYCODE_ENTER } from 'jslib/helpers/keyboard';
import services from '../utils/services';
import BookSelector from './BookSelector';
import Flash from './Flash';
import { useSelector, useDispatch } from '../store/hooks';
import { updateContent, resetComposer } from '../store/composer/actions';
import { fetchBooks } from '../store/books/actions';
import { navigate } from '../store/location/actions';

interface Props {}

// focusBookSelectorInput focuses on the input element of the book selector.
// It needs to traverse the tree returned by the ref API of the 'react-select' library,
// and to guard against possible breaking changes, if the path does not exist, it noops.
function focusBookSelectorInput(bookSelectorRef) {
  bookSelectorRef.select &&
    bookSelectorRef.select.select &&
    bookSelectorRef.select.select.inputRef &&
    bookSelectorRef.select.select.inputRef.focus();
}

function useFetchData() {
  const dispatch = useDispatch();

  const { books } = useSelector(state => {
    return {
      books: state.books
    };
  });

  useEffect(() => {
    if (!books.isFetched) {
      dispatch(fetchBooks());
    }
  }, [dispatch, books.isFetched]);
}

function useInitFocus(contentRef, bookSelectorRef) {
  const { composer, books } = useSelector(state => {
    return {
      composer: state.composer,
      books: state.books
    };
  });

  useEffect(() => {
    if (!books.isFetched) {
      return () => null;
    }

    if (bookSelectorRef && contentRef) {
      if (composer.bookLabel === '') {
        focusBookSelectorInput(bookSelectorRef);
      } else {
        contentRef.focus();
      }
    }
  }, [contentRef, bookSelectorRef, books.isFetched]);
}

const Composer: React.FunctionComponent<Props> = () => {
  useFetchData();
  const [contentFocused, setContentFocused] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [errMsg, setErrMsg] = useState('');
  const dispatch = useDispatch();
  const [contentRef, setContentEl] = useState(null);
  const [bookSelectorRef, setBookSelectorEl] = useState(null);

  const { composer, settings } = useSelector(state => {
    return {
      composer: state.composer,
      settings: state.settings
    };
  });

  const handleSubmit = async e => {
    e.preventDefault();

    setSubmitting(true);

    try {
      let bookUUID;
      if (composer.bookUUID === '') {
        const resp = await services.books.create(
          {
            name: composer.bookLabel
          },
          {
            headers: {
              Authorization: `Bearer ${settings.sessionKey}`
            }
          }
        );

        bookUUID = resp.book.uuid;
      } else {
        bookUUID = composer.bookUUID;
      }

      const resp = await services.notes.create(
        {
          book_uuid: bookUUID,
          content: composer.content
        },
        {
          headers: {
            Authorization: `Bearer ${settings.sessionKey}`
          }
        }
      );

      // clear the composer state
      setErrMsg('');
      setSubmitting(false);

      dispatch(resetComposer());

      // navigate
      dispatch(
        navigate('/success', {
          bookName: composer.bookLabel,
          noteUUID: resp.result.uuid
        })
      );
    } catch (e) {
      setErrMsg(e.message);
      setSubmitting(false);
    }
  };

  const handleSubmitShortcut = e => {
    // Shift + Enter
    if (e.shiftKey && e.keyCode === KEYCODE_ENTER) {
      handleSubmit(e);
    }
  };

  useEffect(() => {
    window.addEventListener('keydown', handleSubmitShortcut);

    return () => {
      window.removeEventListener('keydown', handleSubmitShortcut);
    };
  }, []);

  useEffect(() => {}, []);

  let submitBtnText: string;
  if (submitting) {
    submitBtnText = 'Saving...';
  } else {
    submitBtnText = 'Save';
  }

  useInitFocus(contentRef, bookSelectorRef);

  return (
    <div className="composer">
      <Flash when={errMsg !== ''} message={errMsg} />

      <form onSubmit={handleSubmit} className="form">
        <BookSelector
          selectorRef={setBookSelectorEl}
          onAfterChange={() => {
            contentRef.focus();
          }}
        />

        <div className="content-container">
          <textarea
            className="content"
            placeholder="What did you learn?"
            onChange={e => {
              const val = e.target.value;

              dispatch(updateContent(val));
            }}
            value={composer.content}
            ref={el => {
              setContentEl(el);
            }}
            onFocus={() => {
              setContentFocused(true);
            }}
            onBlur={() => {
              setContentFocused(false);
            }}
          />

          <div
            className={classnames('shortcut-hint', { shown: contentFocused })}
          >
            Shift + Enter to save
          </div>
        </div>

        <input
          type="submit"
          value={submitBtnText}
          className="submit-button"
          disabled={submitting}
        />
      </form>
    </div>
  );
};

export default Composer;
