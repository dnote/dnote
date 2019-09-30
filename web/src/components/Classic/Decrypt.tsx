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
import Helmet from 'react-helmet';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import { b64ToBuf, bufToUtf8 } from 'web/libs/encoding';
import { homePathDef } from 'web/libs/paths';
import services from 'web/libs/services';
import Logo from '../Icons/Logo';
import { aes256GcmDecrypt } from '../../crypto';
import { useDispatch } from '../../store';
import { setMessage } from '../../store/ui';

interface Props extends RouteComponentProps {}

const ClassicDecrypt: React.SFC<Props> = ({ history }) => {
  const [errMsg, setErrMsg] = useState('');
  const [progressMsg, setProgressMsg] = useState('');
  const [busy, setIsBusy] = useState(false);

  const dispatch = useDispatch();

  async function handleDecrypt() {
    try {
      setIsBusy(true);

      const cipherKey = localStorage.getItem('cipherKey');
      const cipherKeyBuf = b64ToBuf(cipherKey);

      const books = await services.books.fetch({ encrypted: true });
      for (let i = 0; i < books.length; i++) {
        const book = books[i];
        const labelBuf = b64ToBuf(book.label);
        console.log('book.label', book);

        console.log('cipherKeyBuf', cipherKeyBuf, 'labelBuf', labelBuf);

        // eslint-disable-next-line no-await-in-loop
        const labelDec = await aes256GcmDecrypt(cipherKeyBuf, labelBuf);

        console.log(labelDec);

        // eslint-disable-next-line no-await-in-loop
        await services.books.update(book.uuid, {
          name: bufToUtf8(labelDec)
        });
      }

      const notes = await services.notes.classicFetch();
      for (let i = 0; i < notes.length; i++) {
        const note = notes[i];

        let contentDec: string;

        if (note.content !== '') {
          const contentBuf = b64ToBuf(note.content);

          if (i % 10 === 0) {
            setProgressMsg(`${i} of ${notes.length} notes decrypted...`);
          }

          // eslint-disable-next-line no-await-in-loop
          const contentDecBuf = await aes256GcmDecrypt(
            cipherKeyBuf,
            contentBuf
          );
          contentDec = bufToUtf8(contentDecBuf);
        } else {
          contentDec = '';
        }

        // eslint-disable-next-line no-await-in-loop
        await services.notes.update(note.uuid, {
          content: contentDec
        });
      }

      await services.users.classicCompleteMigrate();

      dispatch(
        setMessage({
          message:
            'Congratulations. You are now using the new Dnote focusing on knowledge base',
          kind: 'info',
          path: homePathDef
        })
      );

      localStorage.removeItem('cipherKey');

      history.push('/');
    } catch (e) {
      console.log(e);
      setErrMsg(e.message);
      setProgressMsg('');
      setIsBusy(false);
    }
  }

  return (
    <div>
      <Helmet>
        <title>Decrypt (Classic)</title>
      </Helmet>

      <div className="container">
        <a href="/">
          <Logo fill="#252833" width={60} height={60} />
        </a>
        <h1 className="heading">Decrypt your notes and books</h1>

        <div className="auth-body">
          <div className="auth-panel">
            {errMsg && <div className="alert alert-danger">{errMsg}</div>}
            {progressMsg && (
              <div className="alert alert-info">{progressMsg}</div>
            )}

            <p>
              Please press the Decrypt button to decrypt all your notes and
              books.
            </p>

            <button
              onClick={handleDecrypt}
              className="button button-first button-normal"
              type="button"
              disabled={busy}
            >
              {busy ? 'Decrypting...' : 'Decrypt'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default withRouter(ClassicDecrypt);
