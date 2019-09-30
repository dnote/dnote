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
import Helmet from 'react-helmet';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import { legacyMigrate } from 'jslib/services/users';
import { b64ToBuf, utf8ToBuf, bufToB64 } from 'web/libs/encoding';
import * as booksService from 'jslib/services/books';
import * as notesService from 'jslib/services/notes';
import { legacyFetchNotes } from 'jslib/services/notes';
import { updateAuthEmail } from '../../actions/form';
import { receiveUser, legacyGetCurrentUser } from '../../actions/auth';
import Logo from '../Icons/Logo';
import { updateMessage } from '../../actions/ui';
import LegacyFooter from '../Common/LegacyFooter';
import { aes256GcmEncrypt } from '../../crypto';

class LegacyEncrypt extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMsg: '',
      busy: false,
      isReady: false,
      progressMsg: ''
    };
  }

  componentDidMount() {
    const { doLegacyGetCurrentUser } = this.props;

    doLegacyGetCurrentUser().then(() => {
      const { userData, history } = this.props;

      const user = userData.data;

      if (user.encrypted) {
        history.push('/app');
      } else {
        this.setState({ isReady: true });
      }
    });
  }

  handleEncrypt = () => {
    this.setState({ busy: true }, async () => {
      try {
        const cipherKey = localStorage.getItem('cipherKey');
        const cipherKeyBuf = b64ToBuf(cipherKey);

        const books = await booksService.fetch({ encrypted: false });
        for (let i = 0; i < books.length; i++) {
          const book = books[i];
          const labelBuf = utf8ToBuf(book.label);

          // eslint-disable-next-line no-await-in-loop
          const labelEnc = await aes256GcmEncrypt(cipherKeyBuf, labelBuf);

          // eslint-disable-next-line no-await-in-loop
          await booksService.update(book.uuid, {
            name: bufToB64(labelEnc)
          });
        }

        const notes = await legacyFetchNotes({ encrypted: false });
        for (let i = 0; i < notes.length; i++) {
          const note = notes[i];
          const contentBuf = utf8ToBuf(note.content);

          if (i % 10 === 0) {
            this.setState({
              progressMsg: `${i} of ${notes.length} notes encrypted...`
            });
          }

          // eslint-disable-next-line no-await-in-loop
          const contentEnc = await aes256GcmEncrypt(cipherKeyBuf, contentBuf);

          // eslint-disable-next-line no-await-in-loop
          await notesService.update(note.uuid, {
            content: bufToB64(contentEnc)
          });
        }

        await legacyMigrate();

        const { history, doUpdateMessage } = this.props;
        doUpdateMessage(
          'Congratulations. You are now using encrypted Dnote',
          'info'
        );
        history.push('/');
      } catch (e) {
        console.log(e);
        this.setState({ busy: false, errorMsg: e.message, progressMsg: '' });
      }
    });
  };

  render() {
    const { errorMsg, progressMsg, busy, isReady } = this.state;

    if (!isReady) {
      return <div>Loading...</div>;
    }

    return (
      <div className="auth-page page">
        <Helmet>
          <title>Encrypt</title>
        </Helmet>

        <div className="container">
          <div className="container">
            <a href="/">
              <Logo fill="#252833" width="60" height="60" />
            </a>
            <h1 className="heading">Encrypt your notes and books</h1>

            <div className="auth-body">
              <div className="auth-panel">
                {errorMsg && (
                  <div className="alert alert-danger">{errorMsg}</div>
                )}
                {progressMsg && (
                  <div className="alert alert-info">{progressMsg}</div>
                )}

                <p>
                  Please press the Encrypt button to encrypt all your notes and
                  books.
                </p>

                <button
                  onClick={this.handleEncrypt}
                  className="button button-first"
                  type="button"
                  disabled={busy}
                >
                  {busy ? 'Encrypting...' : 'Encrypt'}
                </button>
              </div>
            </div>
          </div>

          <LegacyFooter />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    email: state.form.auth.email,
    userData: state.auth.user
  };
}

const mapDispatchToProps = {
  doUpdateAuthFormEmail: updateAuthEmail,
  doReceiveUser: receiveUser,
  doLegacyGetCurrentUser: legacyGetCurrentUser,
  doUpdateMessage: updateMessage
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(LegacyEncrypt)
);
