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

import settingsStyles from '../Settings.scss';
import styles from './Placeholder.scss';

const Placeholder: React.SFC = () => {
  return (
    <div className="container-wide">
      <div className="row">
        <div className="col-12 col-md-12 col-lg-10">
          <section className={settingsStyles.section}>
            <div className={settingsStyles['section-heading']}>
              <span>&nbsp;</span>
            </div>

            <div className={styles.content1}>
              <div className={classnames('holder', styles['content1-line1'])} />
              <div className={classnames('holder', styles['content1-line2'])} />
              <div className={classnames('holder', styles['content1-line3'])} />
            </div>

            <div className={styles.content2}>
              <div className={styles['content2-left']}>
                <div
                  className={classnames('holder', styles['content2-line1'])}
                />
                <div
                  className={classnames('holder', styles['content2-line2'])}
                />
              </div>

              <div className={styles['content2-right']}>
                <div
                  className={classnames('holder', styles['content2-line2'])}
                />
              </div>
            </div>
          </section>

          <section className={settingsStyles.section}>
            <div className={settingsStyles['section-heading']}>
              <span>&nbsp;</span>
            </div>

            <div className={styles.content3}>
              <div className={styles['content3-left']}>
                <div
                  className={classnames('holder', styles['content3-line1'])}
                />
              </div>

              <div className={styles['content3-right']}>
                <div
                  className={classnames('holder', styles['content3-line2'])}
                />
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
};

export default Placeholder;
