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

/* eslint react/no-danger: 0 */

import React from 'react';
import classnames from 'classnames';

import Flash from '../../Common/Flash';
import { parseMarkdown } from '../../../helpers/markdown';

import styles from './Preview.module.scss';
import workspaceStyle from './Workspace.module.scss';

function Preview({ noteError, content, previewRef }) {
  return (
    <div
      className={classnames(workspaceStyle.pane, styles.wrapper)}
      ref={previewRef}
    >
      {noteError && (
        <Flash type="danger">
          Could not display the note due to an error: {noteError}
        </Flash>
      )}

      <div
        className={classnames(
          'markdown-body',
          styles.content,
          workspaceStyle['pane-content']
        )}
        dangerouslySetInnerHTML={{
          __html: parseMarkdown(content)
        }}
      />
    </div>
  );
}

export default Preview;
