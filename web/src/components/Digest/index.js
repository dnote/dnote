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

import React, { useEffect } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import Helmet from 'react-helmet';

import Note from '../Common/Note';
import { getDigest } from '../../actions/digest';
import { getCipherKey } from '../../crypto';

import styles from './Digest.module.scss';

function Digest({ digestData, doGetDigest, match, demo }) {
  const { digestUUID } = match.params;

  useEffect(() => {
    const cipherKeyBuf = getCipherKey(demo);

    doGetDigest(cipherKeyBuf, digestUUID, demo);
  }, [demo, doGetDigest, digestUUID]);

  const { item } = digestData;

  if (!digestData.isFetched) {
    return null;
  }

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>Digest</title>
      </Helmet>

      <ul className={styles.list}>
        {item.notes.map(note => {
          return (
            <li key={note.uuid} className={styles.item}>
              <Note note={note} />
            </li>
          );
        })}
      </ul>
    </div>
  );
}

function mapStateToProps(state) {
  return {
    digestData: state.digest,
    error: state.digest.error
  };
}

const mapDispatchToProps = {
  doGetDigest: getDigest
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(Digest)
);
