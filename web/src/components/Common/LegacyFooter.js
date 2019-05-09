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

export default () => {
  return (
    <div className="legacy-login-footer">
      <h2>Why do I need this?</h2>
      <p>
        From March 2019, all data on Dnote will be encrypted to respect your
        privacy. To benefit from this new system, you are required to login here
        and set your credentials.
      </p>

      <h2>How long does it take?</h2>
      <p>It will take a minute.</p>

      <h2>What do I need to do?</h2>
      <p>
        After logging in, you will choose an email/password combination. If you
        already have one, you can reuse it. Then, the page will automatically
        encrypt all your data from the beginning of time.
      </p>

      <h2>How will the encryption work?</h2>
      <p>
        Dnote will use AES256 block cipher to encrypt everything before sending
        it to the server. AES256 is considered unbreakable and is used by
        government agencies around the world to encrypt top secrets.
      </p>

      <h2>Can I change email/password later?</h2>
      <p>Yes, you can.</p>
    </div>
  );
};
