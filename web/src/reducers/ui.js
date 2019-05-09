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
  UPDATE_MESSAGE,
  RESET_MESSAGE,
  UPDATE_AUTH_FORM_EMAIL,
  TOGGLE_SIDEBAR,
  CLOSE_SIDEBAR,
  TOGGLE_NOTE_SIDEBAR,
  CLOSE_NOTE_SIDEBAR,
  INIT_LAYOUT
} from '../actions/ui';

const initialState = {
  message: {
    content: '',
    type: ''
  },
  modal: {
    newNote: false
  },
  form: {
    auth: {
      email: ''
    }
  },
  demo: false,
  layout: {
    sidebar: true,
    noteSidebar: true
  }
};

export default function(state = initialState, action) {
  switch (action.type) {
    case UPDATE_MESSAGE: {
      return {
        ...state,
        message: {
          content: action.data.message,
          type: action.data.type
        }
      };
    }
    case RESET_MESSAGE: {
      return {
        ...state,
        message: initialState.message
      };
    }
    case UPDATE_AUTH_FORM_EMAIL: {
      const { data } = action;

      return {
        ...state,
        form: {
          ...state.form,
          auth: {
            ...state.form.auth,
            email: data.email
          }
        }
      };
    }
    case TOGGLE_SIDEBAR: {
      return {
        ...state,
        layout: {
          ...state.layout,
          sidebar: !state.layout.sidebar
        }
      };
    }
    case CLOSE_SIDEBAR: {
      return {
        ...state,
        layout: {
          ...state.layout,
          sidebar: false
        }
      };
    }
    case TOGGLE_NOTE_SIDEBAR: {
      return {
        ...state,
        layout: {
          ...state.layout,
          noteSidebar: !state.layout.noteSidebar
        }
      };
    }
    case CLOSE_NOTE_SIDEBAR: {
      return {
        ...state,
        layout: {
          ...state.layout,
          noteSidebar: false
        }
      };
    }
    case INIT_LAYOUT: {
      return {
        ...state,
        layout: {
          sidebar: action.data.sidebar,
          noteSidebar: action.data.noteSidebar
        }
      };
    }
    default:
      return state;
  }
}
