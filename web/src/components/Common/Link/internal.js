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

function handleLinkClick({ event, editorData }) {
  if (editorData.dirty) {
    const ok = window.confirm('Your unsaved changes will be lost. Continue?');

    if (!ok) {
      event.preventDefault();
    }
  }
}

function mapStateToProps(state) {
  return {
    editorData: state.editor
  };
}

// decorate wraps the given link component from react-rotuer-dom to prevent
// navigation if the current draft is dirty.
export function decorate(LinkComponent, props) {
  function Decorated({ noteData, editorData }) {
    return (
      <LinkComponent
        {...props}
        onClick={e => {
          if (props.onClick) {
            props.onClick(e);
          }

          handleLinkClick({ event: e, noteData, editorData });
        }}
      />
    );
  }

  return connect(mapStateToProps)(Decorated);
}
