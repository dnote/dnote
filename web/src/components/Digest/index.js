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
import { connect } from 'react-redux';
import Helmet from 'react-helmet';

import Note from '../Common/Note';
import { getDigestNotes } from '../../actions/digest';
import { getCipherKey } from '../../crypto';

import './module.scss';

class Digest extends React.Component {
  componentDidMount() {
    const { doGetDigestNotes, match } = this.props;
    const { digestUUID } = match.params;

    // TODO: make demo
    const cipherKeyBuf = getCipherKey();

    doGetDigestNotes(cipherKeyBuf, digestUUID);
  }

  render() {
    const { notes } = this.props;

    return (
      <div className="digest-page">
        <Helmet>
          <title>Digest</title>
        </Helmet>

        <ul className="note-list">
          {notes.map(note => {
            return (
              <li key={note.uuid} className="note-item">
                <Note note={note} />
              </li>
            );
          })}
        </ul>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    notes: state.digest.notes,
    error: state.digest.error
  };
}

const mapDispatchToProps = {
  doGetDigestNotes: getDigestNotes
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Digest);
