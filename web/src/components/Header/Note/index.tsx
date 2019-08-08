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

import React from 'react';
import { useSelector } from '../../../store';

import NormalHeader from '../Normal';
import GuestHeader from './Guest';
import Placeholder from './Placeholder';

interface Props {}

const NoteHeader: React.SFC<Props> = () => {
  const { user } = useSelector(state => {
    return {
      user: state.auth.user
    };
  });

  if (!user.isFetched) {
    return <Placeholder />;
  }

  if (user.data.uuid === '') {
    return <GuestHeader />;
  }

  return <NormalHeader />;
};

export default React.memo(NoteHeader);
