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

import React, { useState, useEffect } from 'react';
import { withRouter } from 'react-router-dom';
import classnames from 'classnames';
import { connect } from 'react-redux';

import { toggleSidebar } from '../../../actions/ui';
import {
  getNotes,
  getMoreNotes,
  getInitialNotes,
  resetNotes
} from '../../../actions/notes';
import { getBooks } from '../../../actions/books';
import { debounce } from '../../../libs/perf';
import NoteGroupList from '../NoteGroupList';
import { parseSearchString } from '../../../libs/url';
import { getScrollbarWidth } from '../../../libs/dom';
import { useEventListener } from '../../../libs/hooks';
import { parsePrevDate } from '../../../libs/notes';
import { getCipherKey } from '../../../crypto';
import BookFilter from './BookFilter';
import { getFacetsFromSearchStr } from '../../../libs/facets';
import { usePrevious } from '../../../libs/hooks';
import { getNotePath, isHomePath } from '../../../libs/paths';
import SidebarToggle from '../../Common/SidebarToggle';
import SubscriberWall from '../../Common/SubscriberWall';
import { isEmptyObj } from '../../../libs/obj';

import styles from './NoteSidebar.module.scss';

let scrollLock;

function useFetchInitialNotes({
  user,
  demo,
  doGetInitialNotes,
  doResetNotes,
  doGetBooks,
  booksData,
  location
}) {
  useEffect(() => {
    if ((isEmptyObj(user) || !user.cloud) && !demo) {
      return () => null;
    }

    const date = new Date();
    const year = date.getUTCFullYear();
    const month = date.getUTCMonth() + 1;
    const facets = getFacetsFromSearchStr(location.search);

    const cipherKeyBuf = getCipherKey(demo);
    doResetNotes();
    doGetInitialNotes({
      facets,
      year,
      month,
      cipherKeyBuf,
      demo
    });

    return () => null;
  }, [
    demo,
    doGetBooks,
    doResetNotes,
    doGetInitialNotes,
    location.search,
    user
  ]);

  useEffect(() => {
    if (!user.cloud || booksData.isFetched) {
      return;
    }

    const cipherKeyBuf = getCipherKey(demo);
    doGetBooks(cipherKeyBuf, demo);
  }, [booksData.isFetched, demo, doGetBooks, user.cloud]);
}

function useRefreshNotes({
  user,
  location,
  demo,
  doGetInitialNotes,
  doResetNotes
}) {
  const prevSearch = usePrevious(location.search);
  const prevIsHome = usePrevious(isHomePath(location.pathname));

  useEffect(() => {
    if (isEmptyObj(user)) {
      return;
    }
    if (prevSearch === null || prevIsHome === null) {
      return;
    }

    const facets = getFacetsFromSearchStr(location.search);
    const isHome = isHomePath(location.pathname);

    if (location.search !== prevSearch || (!prevIsHome && isHome)) {
      const date = new Date();
      const year = date.getUTCFullYear();
      const month = date.getUTCMonth() + 1;

      const cipherKeyBuf = getCipherKey(demo);

      doResetNotes();

      if (user.cloud) {
        doGetInitialNotes({
          cipherKeyBuf,
          facets,
          year,
          month,
          demo
        });
      }
    }
  });
}

function useFetchMoreNotes({
  notesData,
  demo,
  doGetMoreNotes,
  doGetNotes,
  sidebarEl,
  location
}) {
  function fetchMoreNotes() {
    scrollLock = true;

    const queryObj = parseSearchString(location.search);

    if (!notesData.prevDate) {
      return;
    }

    const lastGroup = notesData.groups[notesData.groups.length - 1];

    if (lastGroup.isFetchingMoreNotes) {
      return;
    }

    const cipherKeyBuf = getCipherKey(demo);

    if (lastGroup.total > lastGroup.uuids.length) {
      doGetMoreNotes(
        cipherKeyBuf,
        lastGroup.year,
        lastGroup.month,
        lastGroup.page + 1,
        queryObj,
        demo
      ).then(() => {
        scrollLock = false;
      });
    } else {
      const { year, month } = parsePrevDate(notesData.prevDate);

      doGetNotes(cipherKeyBuf, year, month, queryObj, demo).then(() => {
        scrollLock = false;
      });
    }
  }

  const handleScroll = debounce(() => {
    if (scrollLock || !sidebarEl) {
      return;
    }

    const scrollY = sidebarEl.scrollTop;
    const maxScrollY = sidebarEl.scrollHeight - sidebarEl.clientHeight;

    if (scrollY / maxScrollY > 0.85) {
      fetchMoreNotes();
    }
  }, 100);

  useEventListener(sidebarEl, 'scroll', handleScroll);
}

function usePrevLastGroupItemLength(groups) {
  const lastGroup = groups[groups.length - 1];

  if (!lastGroup) {
    return 0;
  }

  return lastGroup.items.length;
}

function useSelectFirstNote({ notesData, location, history, demo, user }) {
  const { initialized } = notesData;

  const prevDemo = usePrevious(demo);
  const prevInitialized = usePrevious(initialized);
  const prevGroupItemLength = usePrevLastGroupItemLength(notesData.groups);

  useEffect(() => {
    if (prevDemo && !demo && !user.cloud) {
      return;
    }
    if (prevInitialized && prevGroupItemLength > 0) {
      return;
    }

    if (
      !isHomePath(location.pathname, false) &&
      !isHomePath(location.pathname, true)
    ) {
      return;
    }
    if (!initialized) {
      return;
    }

    const firstGroup = notesData.groups[0];
    if (!firstGroup) {
      return;
    }
    if (Object.keys(firstGroup.items).length === 0) {
      return;
    }

    const firstUUID = firstGroup.uuids[0];
    const firstItem = firstGroup.items[firstUUID];
    if (firstItem.errorMessage) {
      return;
    }

    const firstNote = firstItem.data;

    const searchObj = parseSearchString(location.search);
    const dest = getNotePath(firstNote.uuid, searchObj, {
      demo,
      isEditor: true
    });

    history.replace(dest);
  });
}

function useClearDemoData({ demo, doResetNotes }) {
  const prevDemo = usePrevious(demo);

  useEffect(() => {
    if (prevDemo && !demo) {
      doResetNotes();
    }
  });
}

const NoteSidebar = ({
  doToggleSidebar,
  doGetMoreNotes,
  doGetNotes,
  notesData,
  booksData,
  location,
  history,
  match,
  demo,
  layout,
  user,
  doGetInitialNotes,
  doResetNotes,
  doGetBooks
}) => {
  const [sidebarEl, setSidebarEl] = useState(null);
  const [headerEl, setHeaderEl] = useState(null);
  const [bookFilterIsOpen, setBookFilterIsOpen] = useState(false);

  useClearDemoData({ demo, doResetNotes });
  useFetchInitialNotes({
    user,
    demo,
    doGetInitialNotes,
    doResetNotes,
    doGetBooks,
    notesData,
    booksData,
    location
  });
  useRefreshNotes({
    user,
    location,
    demo,
    doGetInitialNotes,
    doResetNotes
  });
  useFetchMoreNotes({
    notesData,
    demo,
    doGetMoreNotes,
    doGetNotes,
    sidebarEl,
    location
  });
  useSelectFirstNote({ notesData, location, history, demo, user });

  useEffect(() => {
    if (!sidebarEl || !headerEl) {
      return;
    }

    const scrollbarWidth = getScrollbarWidth();

    if (bookFilterIsOpen) {
      sidebarEl.style.paddingRight = `${scrollbarWidth}px`;
      headerEl.style.marginRight = `-${scrollbarWidth}px`;
    } else {
      sidebarEl.style.paddingRight = 0;
      headerEl.style.marginRight = 0;
    }
  }, [bookFilterIsOpen, headerEl, sidebarEl]);

  const { noteUUID } = match.params;
  const noteSidebarOpen = layout.noteSidebar;

  return (
    <aside
      className={classnames(styles.sidebar, {
        [styles.hidden]: !noteSidebarOpen,
        [styles.noscroll]: bookFilterIsOpen,
        [styles.fetching]: !notesData.initialized
      })}
      ref={el => {
        setSidebarEl(el);
      }}
    >
      <div className={styles.content}>
        <header
          className={styles.header}
          ref={el => {
            setHeaderEl(el);
          }}
        >
          <div className={styles['header-content']}>
            <div className={styles['header-left']}>
              <SidebarToggle onClick={doToggleSidebar} />

              <div className={styles['header-heading']}>Notes</div>
            </div>

            <div className={styles['header-actions']}>
              <BookFilter
                books={booksData.items}
                isFetching={booksData.isFetching}
                isFetched={booksData.isFetched}
                noteSidebarOpen={layout.noteSidebar}
                isOpen={bookFilterIsOpen}
                setIsOpen={setBookFilterIsOpen}
                demo={demo}
              />
            </div>
          </div>
        </header>

        <NoteGroupList
          groups={notesData.groups}
          demo={demo}
          bookFilterIsOpen={bookFilterIsOpen}
          currentNoteUUID={noteUUID}
          location={location}
          cloud={user.cloud}
        />
        <SubscriberWall
          wrapperClassName={classnames(styles['subscriber-wall'], {
            [styles.pro]: user.cloud
          })}
        />
      </div>
    </aside>
  );
};

function mapStateToProps(state) {
  return {
    notesData: state.notes,
    booksData: state.books,
    layout: state.ui.layout,
    user: state.auth.user.data
  };
}

const mapDispatchToProps = {
  doToggleSidebar: toggleSidebar,
  doGetInitialNotes: getInitialNotes,
  doGetNotes: getNotes,
  doGetMoreNotes: getMoreNotes,
  doGetBooks: getBooks,
  doResetNotes: resetNotes
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(NoteSidebar)
);
