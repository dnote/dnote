/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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
import Helmet from 'react-helmet';
import config from '../../../libs/config';
import { useSelector } from '../../../store';
import SettingRow from '../SettingRow';
import styles from '../Settings.scss';

interface Props {}

const About: React.FunctionComponent<Props> = () => {
  const { user } = useSelector(state => {
    return {
      user: state.auth.user.data
    };
  });

  return (
    <div>
      <Helmet>
        <title>About</title>
      </Helmet>

      <h1 className="sr-only">About</h1>

      <div className={styles.wrapper}>
        <section className={styles.section}>
          <h2 className={styles['section-heading']}>Software</h2>

          <SettingRow name="Version" value={config.version} />
          {!__STANDALONE__ && user.pro && (
            <SettingRow
              name="Support"
              value={<a href="mailto:sung@getdnote.com">sung@getdnote.com</a>}
            />
          )}
        </section>
      </div>
    </div>
  );
};

export default About;
