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
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';

import Helmet from 'react-helmet';
import NoteContent from '../Common/Note';
import Flash from '../Common/Flash';
import { nanosecToMillisec, getShortMonthName } from '../../helpers/time';
import { getNote } from '../../actions/note';
import { getCipherKey } from '../../crypto';
import Placeholder from '../Common/Note/Placeholder';

import styles from './Note.module.scss';

function Note({ demo, match, doGetNote, noteData }) {
  const { params } = match;
  const { noteUUID } = params;

  function formatAddedOnTitle(ts) {
    const ms = nanosecToMillisec(ts);
    const d = new Date(ms);

    const month = getShortMonthName(d);
    const date = d.getDate();
    const year = d.getFullYear();

    return `${month} ${date} ${year}`;
  }

  useEffect(() => {
    const cipherKey = getCipherKey(demo);
    doGetNote(cipherKey, noteUUID, demo);
  }, [demo, doGetNote, noteUUID]);

  const noteError = noteData.error;
  if (noteError) {
    return (
      <Flash type="danger">Error getting the note: {noteData.error}</Flash>
    );
  }

  const { isFetched } = noteData;

  let title;
  if (isFetched) {
    title = `Note: ${formatAddedOnTitle(noteData.item.added_on)}`;
  }

  return (
    <div className={styles.wrapper}>
      <Helmet>
        <title>{title}</title>
      </Helmet>

      <div className={styles.inner}>
        {isFetched ? <NoteContent note={noteData.item} /> : <Placeholder />}
      </div>
    </div>
  );
}

function mapStateToProps(state) {
  return {
    noteData: state.note
  };
}

const mapDispatchToProps = {
  doGetNote: getNote
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(Note)
);
