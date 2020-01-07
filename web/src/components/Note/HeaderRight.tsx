/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import GlobeIcon from '../Icons/Globe';
import Tooltip from '../Common/Tooltip';

interface Props {
  isOwner: boolean;
  isPublic: boolean;
}

const HeaderRight: React.FunctionComponent<Props> = ({ isOwner, isPublic }) => {
  if (!isOwner) {
    return null;
  }
  if (!isPublic) {
    return null;
  }

  const publicTooltip = 'Anyone on the Internet can see this note.';

  return (
    <Tooltip
      id="note-public-indicator"
      alignment="right"
      direction="bottom"
      overlay={publicTooltip}
    >
      <GlobeIcon
        fill="#8c8c8c"
        width={16}
        height={16}
        ariaLabel={publicTooltip}
      />
    </Tooltip>
  );
};

export default HeaderRight;
