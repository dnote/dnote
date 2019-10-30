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

import settingsStyles from '../../Settings.scss';
import styles from './Placeholder.scss';

const Placeholder: React.FunctionComponent = () => {
  return (
    <div className="container-wide">
      <div className="row">
        <div className="col-12 col-md-12 col-lg-10">
          <section className={settingsStyles.section}>
            <div className={styles.content}>
              <div className={styles['content-left']}>
                <div
                  className={classnames('holder', styles['content-line1'])}
                />
              </div>

              <div className={styles['content-right']}>
                <div
                  className={classnames('holder', styles['content-line2'])}
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
