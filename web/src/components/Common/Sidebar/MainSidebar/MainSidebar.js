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

import React, { useCallback, useRef, useEffect } from 'react';
import classnames from 'classnames';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';

import SafeNavLink from '../../Link/SafeNavLink';
import SafeLink from '../../Link/SafeLink';

import {
  getHomePath,
  getBooksPath,
  getNotePath,
  getDigestsPath,
  getSubscriptionPath
} from 'web/libs/paths';
import { parseSearchString } from 'jslib/helpers/url';
import {
  getWindowWidth,
  noteSidebarThreshold,
  sidebarOverlayThreshold
} from 'jslib/helpers/ui';
import { addNote } from '../../../../actions/notes';
import { addBook } from '../../../../actions/books';
import { closeSidebar, closeNoteSidebar } from '../../../../actions/ui';
import { getFacetsFromSearchStr } from 'jslib/helpers/facets';
import NoteIcon from '../../../Icons/Note';
import CloseIcon from '../../../Icons/Close';
import BookIcon from '../../../Icons/Book';
import DashboardIcon from '../../../Icons/Dashboard';
import * as notesOperation from 'jslib/operations/notes';
import * as booksOperation from 'jslib/operations/books';

import styles from './MainSidebar.module.scss';
import sidebarStyles from '../Sidebar.module.scss';

async function getTargetBookUUID({ location, books }) {
  if (books.length === 0) {
    return null;
  }

  const facets = getFacetsFromSearchStr(location.search);
  if (facets.book) {
    return facets.book;
  }

  return books[0].uuid;
}

async function createDefaultBook({ doAddBook }) {
  try {
    const book = await booksOperation.create({ name: 'Default' });
    doAddBook(book);

    return book;
  } catch (err) {
    throw new Error(`creating a book: ${err.message}`);
  }
}

async function getTargetBook({ location, books, doAddBook }) {
  const bookUUID = await getTargetBookUUID({ location, books });
  if (!bookUUID) {
    return createDefaultBook({ doAddBook });
  }

  try {
    return booksOperation.get(bookUUID);
  } catch (err) {
    if (err.response.status === 404) {
      return createDefaultBook({ doAddBook });
    }

    throw new Error(`getting the book: ${err.message}`);
  }
}

async function handleCreateNote({
  books,
  doAddNote,
  doAddBook,
  doCloseSidebar,
  doCloseNoteSidebar,
  history,
  location
}) {
  const width = getWindowWidth();
  if (width < noteSidebarThreshold) {
    doCloseSidebar();
    doCloseNoteSidebar();
  }

  let book;
  try {
    book = await getTargetBook({
      location,
      books,
      doAddBook
    });
  } catch (err) {
    throw new Error(`Checking if book exists ${err.message}`);
  }

  try {
    const note = await notesOperation.create({
      bookUUID: book.uuid,
      content: ''
    });
    const date = new Date();
    const year = date.getUTCFullYear();
    const month = date.getUTCMonth() + 1;

    doAddNote(note, year, month);

    const queryObj = parseSearchString(location.search);
    const dest = getNotePath(note.uuid, queryObj, { isEditor: true });
    history.push(dest);
  } catch (e) {
    console.log('err', e);
  }
}

const Sidebar = ({
  demo,
  location,
  layoutData,
  booksData,
  userData,
  doAddNote,
  doAddBook,
  history,
  doCloseSidebar,
  doCloseNoteSidebar
}) => {
  const sidebarRef = useRef(null);

  const handleLinkClick = useCallback(() => {
    const width = getWindowWidth();

    if (width < noteSidebarThreshold) {
      doCloseSidebar();
    }
  }, [doCloseSidebar]);

  useEffect(() => {
    function handleMousedown(e) {
      const sidebarEl = sidebarRef.current;

      if (sidebarEl && !sidebarEl.contains(e.target)) {
        doCloseSidebar();
      }
    }

    const width = getWindowWidth();
    if (layoutData.sidebar && width < sidebarOverlayThreshold) {
      document.addEventListener('mousedown', handleMousedown);

      return () => {
        document.removeEventListener('mousedown', handleMousedown);
      };
    }

    return () => null;
  }, [layoutData.sidebar, doCloseSidebar]);

  const pathHome = getHomePath({}, { demo });
  const pathBooks = getBooksPath({ demo });
  const pathDigests = getDigestsPath({ demo });

  const user = userData.data;

  return (
    <aside
      className={classnames(sidebarStyles.wrapper, {
        [sidebarStyles['wrapper-hidden']]: !layoutData.sidebar
      })}
    >
      <button
        aria-label="Close the sidebar"
        type="button"
        className={classnames(sidebarStyles['close-button'], {
          [sidebarStyles['close-button-hidden']]: !layoutData.sidebar
        })}
        onClick={e => {
          e.preventDefault();
          doCloseSidebar();
        }}
      >
        <div className={sidebarStyles['close-button-content']}>
          <CloseIcon width={16} height={16} />
        </div>
      </button>
      <div
        className={classnames(sidebarStyles.sidebar, {
          [sidebarStyles['sidebar-hidden']]: !layoutData.sidebar
        })}
        ref={sidebarRef}
      >
        <div className={classnames(sidebarStyles['sidebar-content'])}>
          <div>
            <div className={styles['button-wrapper']}>
              <button
                id="T-create-note-btn"
                type="button"
                className={classnames(
                  'button button-normal button-slim button-stretch button-third'
                )}
                onClick={() => {
                  handleCreateNote({
                    books: booksData.items,
                    doAddNote,
                    doAddBook,
                    doCloseSidebar,
                    doCloseNoteSidebar,
                    history,
                    location
                  });
                }}
                disabled={demo || !userData.isFetched || !user.cloud}
              >
                New note
              </button>
            </div>

            <ul
              className={classnames(
                'list-unstyled',
                sidebarStyles['link-list']
              )}
            >
              <li className={sidebarStyles['link-item']}>
                <SafeNavLink
                  onClick={handleLinkClick}
                  className={sidebarStyles.link}
                  to={pathHome}
                  activeClassName={sidebarStyles['link-active']}
                  isActive={() => {
                    return location.pathname === pathHome.pathname;
                  }}
                >
                  <NoteIcon width={16} height={16} fill="#6e6e6e" />
                  <div className={sidebarStyles['link-label']}>All notes</div>
                </SafeNavLink>
              </li>

              <li className={sidebarStyles['link-item']}>
                <SafeNavLink
                  onClick={handleLinkClick}
                  className={sidebarStyles.link}
                  to={pathBooks}
                  activeClassName={sidebarStyles['link-active']}
                  isActive={() => {
                    return location.pathname === pathBooks.pathname;
                  }}
                >
                  <BookIcon width={16} height={16} fill="#6e6e6e" />
                  <div className={sidebarStyles['link-label']}>Books</div>
                </SafeNavLink>
              </li>

              <li className={sidebarStyles['link-item']}>
                <SafeNavLink
                  onClick={handleLinkClick}
                  className={sidebarStyles.link}
                  to={pathDigests}
                  activeClassName={sidebarStyles['link-active']}
                  isActive={() => {
                    return location.pathname === pathDigests.pathname;
                  }}
                >
                  <DashboardIcon width={16} height={16} fill="#6e6e6e" />
                  <div className={sidebarStyles['link-label']}>Digests</div>
                </SafeNavLink>
              </li>
            </ul>

            <div
              className={classnames(styles['upgrade-wrapper'], {
                [styles['upgrade-wrapper-shown']]: !user.cloud
              })}
            >
              <SafeLink
                to={getSubscriptionPath()}
                onClick={handleLinkClick}
                className="button button-slim button-stretch button-first"
              >
                Upgrade
              </SafeLink>
            </div>
          </div>

          {/*
          <div className={styles.bottom}>
            <span className={styles.version}>v1.0.0</span>
          </div>
            */}
        </div>
      </div>
    </aside>
  );
};

function mapStateToProps(state) {
  return {
    layoutData: state.ui.layout,
    booksData: state.books,
    userData: state.auth.user
  };
}

const mapDispatchToProps = {
  doAddBook: addBook,
  doAddNote: addNote,
  doCloseSidebar: closeSidebar,
  doCloseNoteSidebar: closeNoteSidebar
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(Sidebar)
);
