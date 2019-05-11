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
import { connect } from 'react-redux';
import Helmet from 'react-helmet';

import { getDigests, getMoreDigests } from '../../actions/digests';
import { debounce } from '../../libs/perf';
import { useEventListener } from '../../libs/hooks';

import Header from '../Common/Page/Header';
import Body from '../Common/Page/Body';
import SubscriberWall from '../Common/SubscriberWall';
import Content from './Content';

function useFetchMoreDigests({ digestsData, pageEl, doGetMoreDigests }) {
  let scrollLock = false;

  function fetchMore() {
    scrollLock = true;

    doGetMoreDigests().then(() => {
      scrollLock = false;
    });
  }

  const handleScroll = debounce(() => {
    if (scrollLock || !pageEl) {
      return;
    }

    const scrollY = pageEl.scrollTop;
    const maxScrollY = pageEl.scrollHeight - pageEl.clientHeight;

    if (scrollY / maxScrollY > 0.85) {
      if (digestsData.total > digestsData.items.length) {
        fetchMore();
      }
    }
  }, 100);

  useEventListener(pageEl, 'scroll', handleScroll);
}

function Digests({
  demo,
  userData,
  digestsData,
  doGetDigests,
  doGetMoreDigests
}) {
  const [pageEl, setPageEl] = useState(null);

  useEffect(() => {
    doGetDigests({ demo });
  }, [demo, doGetDigests]);

  const user = userData.data;

  useFetchMoreDigests({ digestsData, pageEl, doGetMoreDigests });

  return (
    <div
      className="page"
      ref={el => {
        setPageEl(el);
      }}
    >
      <Helmet>
        <title>Digests</title>
      </Helmet>

      <Header heading="Digests" />

      <Body>
        <div className="container">
          {demo || user.cloud ? (
            <Content user={user} demo={demo} />
          ) : (
            <SubscriberWall />
          )}
        </div>
      </Body>
    </div>
  );
}

function mapStateToProps(state) {
  return {
    userData: state.auth.user,
    digestsData: state.digests
  };
}

const mapDispatchToProps = {
  doGetDigests: getDigests,
  doGetMoreDigests: getMoreDigests
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Digests);
