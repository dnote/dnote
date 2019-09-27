import { useEffect } from 'react';

import { focusTextarea } from 'web/libs/dom';
import { resetEditor } from '../../store/editor';
import { useDispatch } from '../../store';

// useFocusTextarea is a hook that, when the given textareaEl becomes
// defined, focuses on it.
export function useFocusTextarea(textareaEl: HTMLTextAreaElement) {
  useEffect(() => {
    if (textareaEl) {
      focusTextarea(textareaEl);
    }
  }, [textareaEl]);
}

// useCleanupEditor is a hook that cleans up the editor state
export function useCleanupEditor() {
  const dispatch = useDispatch();

  useEffect(() => {
    return () => {
      dispatch(resetEditor());
    };
  }, [dispatch]);
}
