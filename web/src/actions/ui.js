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

export const UPDATE_MESSAGE = 'ui/UPDATE_MESSAGE';
export const RESET_MESSAGE = 'ui/RESET_MESSAGE';
export const UPDATE_AUTH_FORM_EMAIL = 'ui/UPDATE_AUTH_FORM_EMAIL';
export const TOGGLE_SIDEBAR = 'ui/TOGGLE_SIDEBAR';
export const CLOSE_SIDEBAR = 'ui/CLOSE_SIDEBAR';
export const TOGGLE_NOTE_SIDEBAR = 'ui/TOGGLE_NOTE_SIDEBAR';
export const CLOSE_NOTE_SIDEBAR = 'ui/CLOSE_NOTE_SIDEBAR';
export const INIT_LAYOUT = 'ui/INIT_LAYOUT';

export function updateMessage(message, type) {
  return {
    type: UPDATE_MESSAGE,
    data: { message, type }
  };
}

export function resetMessage() {
  return {
    type: RESET_MESSAGE
  };
}

export function updateAuthFormEmail(email) {
  return {
    type: UPDATE_AUTH_FORM_EMAIL,
    data: {
      email
    }
  };
}

export function toggleSidebar() {
  return {
    type: TOGGLE_SIDEBAR
  };
}

export function toggleNoteSidebar() {
  return {
    type: TOGGLE_NOTE_SIDEBAR
  };
}

export function closeNoteSidebar() {
  return {
    type: CLOSE_NOTE_SIDEBAR
  };
}

export function closeSidebar() {
  return {
    type: CLOSE_SIDEBAR
  };
}

export function initLayout(data) {
  return {
    type: INIT_LAYOUT,
    data
  };
}
