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

import {
  ADD,
  REFRESH,
  RECEIVE,
  START_FETCHING,
  RECEIVE_MORE,
  START_FETCHING_MORE,
  RECEIVE_ERROR,
  RESET,
  REMOVE,
  NotesActionType,
  NotesState,
  NotesGroup,
  NotesGroupItems
} from './type';
import { NoteData } from '../../operations/types';
import { removeKey } from '../../libs/obj';

function now(): number {
  const d = new Date();
  const ts = d.getTime();

  return ts * 1000000;
}

const initialState: NotesState = {
  groups: [],
  initialized: false,
  prevDate: now()
};

function makeGroup(year, month): NotesGroup {
  return {
    year,
    month,
    uuids: [],
    items: {},
    total: 0,
    isFetching: true,
    isFetched: false,
    isFetchingMore: false,
    hasFetchedMore: false,
    page: 0,
    error: null
  };
}

function toNoteItems(notes: NoteData[]): NotesGroupItems {
  const ret = {};

  for (let i = 0; i < notes.length; ++i) {
    const note = notes[i];

    ret[note.uuid] = note;
  }

  return ret;
}

export default function(
  state = initialState,
  action: NotesActionType
): NotesState {
  switch (action.type) {
    case START_FETCHING: {
      const { groups } = state;
      const { year, month } = action.data;

      // If duplicate exists, keep original state
      for (let i = 0; i < groups.length; i++) {
        const group = groups[i];

        if (group.year === year && group.month === month) {
          return state;
        }
      }

      return {
        ...state,
        groups: [...groups, makeGroup(year, month)]
      };
    }
    case ADD: {
      const { note, year, month } = action.data;

      let exists = false;
      let groups = state.groups.map(group => {
        if (group.year === year && group.month === month) {
          exists = true;

          return {
            ...group,
            uuids: [note.uuid, ...group.uuids],
            items: { ...toNoteItems([note]), ...group.items },
            total: group.total + 1
          };
        }

        return group;
      });

      // If this is the first entry in the group, make and prepend a new group
      // with this note
      if (!exists) {
        groups = [
          {
            ...makeGroup(year, month),
            items: toNoteItems([note]),
            uuids: [note.uuid],
            total: 1,
            isFetching: false,
            isFetched: true,
            isFetchingMore: false,
            hasFetchedMore: false,
            page: 1
          },
          ...state.groups
        ];
      }

      return {
        ...state,
        groups
      };
    }
    case REFRESH: {
      const { year, month, noteUUID, book, content, isPublic } = action.data;

      const groups = state.groups.map(group => {
        if (group.year === year && group.month === month) {
          return {
            ...group,
            items: {
              ...group.items,
              [noteUUID]: {
                ...group.items[noteUUID],
                data: {
                  ...group.items[noteUUID],
                  book,
                  content,
                  public: isPublic
                }
              }
            }
          };
        }

        return group;
      });

      return {
        ...state,
        groups
      };
    }
    case REMOVE: {
      const { year, month, noteUUID } = action.data;

      const groups = [];

      state.groups.forEach(group => {
        if (group.year === year && group.month === month) {
          // If this item is the last of the group, dismiss the group
          if (group.total > 1) {
            const filterItems = removeKey(group.items, noteUUID);

            const g = {
              ...group,
              uuids: group.uuids.filter(uuid => {
                return uuid !== noteUUID;
              }),
              items: filterItems,
              total: group.total - 1
            };

            groups.push(g);
          }
        } else {
          groups.push(group);
        }
      });

      return {
        ...state,
        groups
      };
    }
    case RECEIVE: {
      const { notes, year, month, total, prevDate } = action.data;

      let groups;
      if (notes.length > 0) {
        groups = state.groups.map(group => {
          if (group.year === year && group.month === month) {
            const uuids = notes.map(note => {
              return note.uuid;
            });

            return {
              ...group,
              uuids,
              items: toNoteItems(notes),
              isFetching: false,
              isFetched: true,
              page: 1,
              total
            };
          }

          return group;
        });
      } else {
        groups = state.groups.filter(group => {
          return group.year !== year || group.month !== month;
        });
      }

      return {
        ...state,
        prevDate,
        groups,
        initialized: true
      };
    }
    case START_FETCHING_MORE: {
      const { year, month } = action.data;

      return {
        ...state,
        groups: state.groups.map(group => {
          if (group.year === year && group.month === month) {
            return {
              ...group,
              isFetchingMore: true,
              hasFetchedMore: false
            };
          }

          return group;
        })
      };
    }
    case RECEIVE_MORE: {
      const { year, month, notes, prevDate } = action.data;

      return {
        ...state,
        prevDate,
        groups: state.groups.map(group => {
          if (group.year === year && group.month === month) {
            const uuids = notes.map(note => {
              return note.uuid;
            });

            return {
              ...group,
              uuids: [...uuids, ...group.uuids],
              items: { ...toNoteItems(notes), ...group.items },
              page: group.page + 1,
              isFetchingMore: false,
              hasFetchedMore: true
            };
          }

          return group;
        })
      };
    }
    case RECEIVE_ERROR: {
      return {
        ...state,
        groups: state.groups.map(group => {
          const { year, month, error } = action.data;

          if (group.year === year && group.month === month) {
            return {
              ...group,
              error,
              isFetchingMore: false,
              hasFetchedMore: false,
              isFetching: false,
              isFetched: false
            };
          }

          return group;
        })
      };
    }
    case RESET: {
      return initialState;
    }
    default:
      return state;
  }
}
