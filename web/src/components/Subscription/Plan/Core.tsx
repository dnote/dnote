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

import { UserData } from 'jslib/operations/types';
import Plan from './internal';

const desc =
  'Streamline your learnings into a personal knowledge base. You can access any item at any time.';

interface Props {
  wrapperClassName: string;
  user: UserData;
  bottomContent: React.ReactElement;
}

const Core: React.FunctionComponent<Props> = ({
  wrapperClassName,
  user,
  bottomContent
}) => {
  return (
    <Plan
      name="Core"
      desc={desc}
      price="Free"
      wrapperClassName={wrapperClassName}
      ctaContent={
        <button
          type="button"
          className="button button-large button-second button-stretch"
          disabled
        >
          {user && user.pro ? 'Already upgraded!' : 'Your current plan'}
        </button>
      }
      bottomContent={bottomContent}
    />
  );
};

export default Core;
