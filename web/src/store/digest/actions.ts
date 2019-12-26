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

import operations from 'web/libs/operations';
import services from 'web/libs/services';
import { DigestData } from 'jslib/operations/types';
import {
  RECEIVE,
  START_FETCHING,
  ERROR,
  RESET,
  SET_NOTE_REVIEWED,
  SetNoteReviewed
} from './type';
import { ThunkAction } from '../types';

export function receiveDigest(digest: DigestData) {
  return {
    type: RECEIVE,
    data: { digest }
  };
}

export function resetDigest() {
  return {
    type: RESET
  };
}

function startFetchingDigest() {
  return {
    type: START_FETCHING
  };
}

function receiveDigestError(errorMessage: string) {
  return {
    type: ERROR,
    data: { errorMessage }
  };
}

function setNoteReviewed(
  noteUUID: string,
  isReviewed: boolean
): SetNoteReviewed {
  return {
    type: SET_NOTE_REVIEWED,
    data: {
      noteUUID,
      isReviewed
    }
  };
}

interface GetDigestFacets {
  q?: string;
}

export const getDigest = (
  digestUUID: string
): ThunkAction<DigestData | void> => {
  return dispatch => {
    dispatch(startFetchingDigest());

    return operations.digests
      .fetch(digestUUID)
      .then(digest => {
        dispatch(receiveDigest(digest));

        return digest;
      })
      .catch(err => {
        console.log('getDigest error', err.message);
        dispatch(receiveDigestError(err.message));
      });
  };
};

export const setDigestNoteReviewed = ({
  digestUUID,
  noteUUID,
  isReviewed
}: {
  digestUUID: string;
  noteUUID: string;
  isReviewed: boolean;
}): ThunkAction<void> => {
  return dispatch => {
    if (!isReviewed) {
      return services.noteReviews.remove(noteUUID).then(() => {
        dispatch(setNoteReviewed(noteUUID, false));
      });
    }

    return services.noteReviews.create({ digestUUID, noteUUID }).then(() => {
      dispatch(setNoteReviewed(noteUUID, true));
    });
  };
};
