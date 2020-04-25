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

import { UPDATE, RESET, SettingsState, SettingsActionType } from './types';
import config from '../../utils/config';

const initialState: SettingsState = {
  apiUrl: config.defaultApiEndpoint,
  webUrl: config.defaultWebUrl
};

export default function (
  state = initialState,
  action: SettingsActionType
): SettingsState {
  switch (action.type) {
    case UPDATE:
      return {
        ...state,
        ...action.data.settings
      };
    case RESET:
      return initialState;
    default:
      return state;
  }
}
