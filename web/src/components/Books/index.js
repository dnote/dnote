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

import React, { useEffect, useState } from 'react';
import Helmet from 'react-helmet';
import { connect } from 'react-redux';
import classnames from 'classnames';

import { getBooks } from '../../actions/books';
import { getCipherKey } from '../../crypto';
import BookPlusIcon from '../Icons/BookPlus';
import Tooltip from '../Common/Tooltip';

import SubscriberWall from '../Common/SubscriberWall';
import Header from '../Common/Page/Header';
import Body from '../Common/Page/Body';
import Content from './Content';
import CreateBookModal from './CreateBookModal';

import styles from './Books.module.scss';

function ContentWrapper({
  userData,
  booksData,
  demo,
  containerEl,
  onStartCreateBook
}) {
  const user = userData.data;

  if (demo || user.cloud) {
    return (
      <Content
        books={booksData.items}
        isFetching={booksData.isFetching}
        isFetched={booksData.isFetched}
        demo={demo}
        containerEl={containerEl}
        onStartCreateBook={onStartCreateBook}
      />
    );
  }

  return <SubscriberWall />;
}

function Total({ isFetching, count }) {
  return (
    <div className={classnames(styles.count, { [styles.hidden]: isFetching })}>
      {count} total
    </div>
  );
}

function Books({ demo, doGetBooks, userData, booksData }) {
  const [containerEl, setContainerEl] = useState(null);
  const [isCreateBookModalOpen, setIsCreateBookModalOpen] = useState(false);

  useEffect(() => {
    const cipherKeyBuf = getCipherKey(demo);
    doGetBooks(cipherKeyBuf, demo);
  }, [doGetBooks, demo]);

  return (
    <div
      className="page"
      ref={el => {
        setContainerEl(el);
      }}
    >
      <Helmet>
        <title>Books</title>
      </Helmet>

      <Header
        heading="Books"
        leftContent={
          <Total
            isFetching={booksData.isFetching}
            count={booksData.items.length}
          />
        }
        rightContent={
          <Tooltip
            id="tooltip-new-book"
            alignment="right"
            direction="bottom"
            overlay={<span>Create a new book</span>}
          >
            <button
              type="button"
              aria-label="Create a new book"
              className={classnames('button-no-ui', styles['new-book-button'])}
              onClick={() => {
                setIsCreateBookModalOpen(!isCreateBookModalOpen);
              }}
              disabled={booksData.isFetching}
            >
              <BookPlusIcon width={24} height={24} fill="black" />
            </button>
          </Tooltip>
        }
      />

      <Body>
        <div className="container-wide">
          <div className="row">
            <div className="col-12">
              <ContentWrapper
                userData={userData}
                booksData={booksData}
                demo={demo}
                containerEl={containerEl}
                onStartCreateBook={setIsCreateBookModalOpen}
              />
            </div>
          </div>
        </div>
      </Body>

      <CreateBookModal
        isOpen={isCreateBookModalOpen}
        onDismiss={() => {
          setIsCreateBookModalOpen(false);
        }}
        demo={demo}
      />
    </div>
  );
}

function mapStateToProps(state) {
  return {
    userData: state.auth.user,
    booksData: state.books
  };
}

const mapDispatchToProps = {
  doGetBooks: getBooks
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Books);
