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

import BookSelector from './BookSelector';
import NoteActions from './NoteActions';
import SidebarToggle from '../../Common/SidebarToggle';
import SubscriberWall from '../../Common/SubscriberWall';
import { toggleNoteSidebar } from '../../../actions/ui';
import Status from './Status';
import Body from './Body';

import styles from './Workspace.module.scss';

const Workspace = ({
  editorData,
  textareaRef,
  noteData,
  doToggleNoteSidebar,
  demo,
  userData
}) => {
  const isReady = noteData.isFetched && !noteData.error;

  const user = userData.data;

  return (
    <div className={styles.wrapper}>
      <div className={styles.actions}>
        <div className={styles['action-top']}>
          <div className={styles['action-left']}>
            <SidebarToggle type="arrow" onClick={doToggleNoteSidebar} />
            <div className={styles['action-heading']}>Note</div>
          </div>

          <div className={styles['action-right']}>
            <BookSelector
              wrapperClassName={styles['desktop-book-selector']}
              isReady={isReady}
              demo={demo}
            />
            <div className={styles['action-right-right']}>
              <Status />
              <NoteActions
                noteUUID={editorData.noteUUID}
                disabled={!isReady}
                demo={demo}
              />
            </div>
          </div>
        </div>

        <div className={styles['action-bottom']}>
          <BookSelector isReady={isReady} demo={demo} />
        </div>
      </div>

      <div className={styles.main}>
        {demo || user.cloud ? (
          <Body
            isReady={isReady}
            editorData={editorData}
            textareaRef={textareaRef}
            noteData={noteData}
            demo={demo}
          />
        ) : (
          <SubscriberWall />
        )}
      </div>
    </div>
  );
};

function mapStateToProps(state) {
  return {
    editorData: state.editor,
    noteData: state.note,
    userData: state.auth.user
  };
}

const mapDispatchToProps = {
  doToggleNoteSidebar: toggleNoteSidebar
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Workspace);
