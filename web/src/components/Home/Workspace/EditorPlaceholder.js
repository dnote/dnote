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
import classnames from 'classnames';

import styles from './EditorPlaceholder.module.scss';
import workspaceStyle from './Workspace.module.scss';

function EditorPlaceholder() {
  return (
    <div className={classnames(workspaceStyle['editor-area'])}>
      <div className={classnames('holder', styles.line, styles.line1)} />
      <div className={classnames('holder', styles.line, styles.line2)} />
      <div className={classnames('holder', styles.line, styles.line3)} />
      <div className={classnames('holder', styles.line, styles.line4)} />
      <div className={classnames('holder', styles.line, styles.line5)} />
      <div className={classnames('holder', styles.linebreak)} />
      <div className={classnames('holder', styles.line, styles.line6)} />
      <div className={classnames('holder', styles.line, styles.line7)} />
      <div className={classnames('holder', styles.line, styles.line8)} />
      <div className={classnames('holder', styles.line, styles.line9)} />
    </div>
  );
}

export default EditorPlaceholder;
