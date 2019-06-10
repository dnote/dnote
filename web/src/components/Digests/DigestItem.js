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

import React, { useState } from 'react';
import classnames from 'classnames';
import moment from 'moment';

import { Link } from 'react-router-dom';

import styles from './DigestItem.module.scss';
import { getDigestPath } from '../../libs/paths';

function DigestItem({ digest, demo }) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <li
      className={classnames(styles.wrapper, {
        [styles.active]: isHovered
      })}
      key={digest.uuid}
      onMouseEnter={() => {
        setIsHovered(true);
      }}
      onMouseLeave={() => {
        setIsHovered(false);
      }}
    >
      <Link className={styles.link} to={getDigestPath(digest.uuid, { demo })}>
        {moment(digest.created_at).format('YYYY MMM Do')}
      </Link>
    </li>
  );
}

export default DigestItem;
