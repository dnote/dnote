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

import { useRef, useEffect, useState } from 'react';

// usePrevious is a hook that saves the current value to be used later
export function usePrevious<T>(value: T): T | null {
  const ref = useRef<T | null>(null);
  useEffect(() => {
    ref.current = value;
  });

  return ref.current;
}

// useEventListener adds an event listener to the target and clears it when the
// component unmounts
export function useEventListener(target, type, listener) {
  const noop = e => e;
  const savedListener = useRef(noop);

  useEffect(() => {
    savedListener.current = listener;
  });

  useEffect(() => {
    if (!target) {
      return () => null;
    }

    function fn(e) {
      savedListener.current(e);
    }

    target.addEventListener(type, fn);

    return () => {
      target.removeEventListener(type, fn);
    };
  }, [type, target]);
}

// useScript loads a third party script
export function useScript(src: string): [boolean, string] {
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (loaded || error) {
      return () => null;
    }

    const script = document.createElement('script');
    script.src = src;
    script.async = true;

    function onLoad() {
      setLoaded(true);
    }
    function onError(err) {
      setError(err.message);
      script.remove();
    }

    script.addEventListener('load', onLoad);
    script.addEventListener('error', onError);

    document.head.appendChild(script);

    return () => {
      script.removeEventListener('load', onLoad);
      script.removeEventListener('error', onError);
    };
  }, [src, loaded, error]);

  return [loaded, error];
}
