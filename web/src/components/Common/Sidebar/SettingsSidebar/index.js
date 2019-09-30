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
import { withRouter, Link, NavLink } from 'react-router-dom';
import { connect } from 'react-redux';

import { getHomePath, getSettingsPath } from 'web/libs/paths';
import {
  getWindowWidth,
  noteSidebarThreshold,
  sidebarOverlayThreshold
} from 'jslib/helpers/ui';
import { closeSidebar, closeNoteSidebar } from '../../../../actions/ui';
import ArrowIcon from '../../../Icons/Arrow';
import CloseIcon from '../../../Icons/Close';

import styles from './SettingsSidebar.module.scss';
import sidebarStyles from '../Sidebar.module.scss';

const SettingsSidebar = ({ location, layoutData, doCloseSidebar }) => {
  const sidebarRef = useRef(null);

  const maybeCloseSidebar = useCallback(() => {
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
          // [styles['close-button-hidden']]: !layoutData.sidebar
        })}
        onClick={maybeCloseSidebar}
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
              <Link
                to={getHomePath()}
                onClick={maybeCloseSidebar}
                className={classnames(
                  'button button-normal button-slim button-stretch button-third-outline',
                  styles['back-button']
                )}
              >
                <ArrowIcon width="16" height="16" fill="#4d4d8b" />
                <span className={styles['button-text']}>Back to notes</span>
              </Link>
            </div>

            <strong className={styles['sidebar-heading']}>Settings</strong>

            <ul
              className={classnames(
                'list-unstyled',
                sidebarStyles['link-list'],
                styles['link-list']
              )}
            >
              <li className={classnames(sidebarStyles['link-item'])}>
                <NavLink
                  onClick={maybeCloseSidebar}
                  className={classnames(styles.link, sidebarStyles.link)}
                  to={getSettingsPath('account')}
                  activeClassName={classnames(
                    sidebarStyles['link-active'],
                    styles['link-active']
                  )}
                  isActive={() => {
                    return location.pathname === getSettingsPath('account');
                  }}
                >
                  {/* <UserIcon width="16" height="16" fill="#6e6e6e" />
                   */}
                  <div className={sidebarStyles['link-label']}>Account</div>
                </NavLink>
              </li>

              <li className={sidebarStyles['link-item']}>
                <NavLink
                  onClick={maybeCloseSidebar}
                  className={classnames(styles.link, sidebarStyles.link)}
                  to={getSettingsPath('notification')}
                  activeClassName={classnames(
                    sidebarStyles['link-active'],
                    styles['link-active']
                  )}
                  isActive={() => {
                    return (
                      location.pathname === getSettingsPath('notification')
                    );
                  }}
                >
                  {/* <EmailIcon width="16" height="16" fill="#6e6e6e" />
                   */}
                  <div className={sidebarStyles['link-label']}>
                    Notification
                  </div>
                </NavLink>
              </li>

              <li className={sidebarStyles['link-item']}>
                <NavLink
                  onClick={maybeCloseSidebar}
                  className={classnames(styles.link, sidebarStyles.link)}
                  to={getSettingsPath('billing')}
                  activeClassName={classnames(
                    sidebarStyles['link-active'],
                    styles['link-active']
                  )}
                  isActive={() => {
                    return location.pathname === getSettingsPath('billing');
                  }}
                >
                  {/* <CreditCardIcon width="16" height="16" fill="#6e6e6e" />
                   */}
                  <div className={sidebarStyles['link-label']}>Billing</div>
                </NavLink>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </aside>
  );
};

function mapStateToProps(state) {
  return {
    layoutData: state.ui.layout,
    booksData: state.books
  };
}

const mapDispatchToProps = {
  doCloseSidebar: closeSidebar,
  doCloseNoteSidebar: closeNoteSidebar
};

export default withRouter(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(SettingsSidebar)
);
