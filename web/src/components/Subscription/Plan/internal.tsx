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

import styles from './Plan.scss';

interface Perk {
  id: string;
  icon: JSX.Element;
  value: string;
}

interface Props {
  name: string;
  desc: string;
  price: string;
  bottomContent: React.ReactElement;
  ctaContent: React.ReactElement;
  interval?: string;
  wrapperClassName: string;
}

const Plan: React.FunctionComponent<Props> = ({
  name,
  desc,
  price,
  bottomContent,
  ctaContent,
  interval,
  wrapperClassName
}) => {
  return (
    <div className={classnames(styles.wrapper, wrapperClassName)}>
      <div
        className={classnames(styles.header, {
          [styles.pro]: name === 'Pro'
        })}
      >
        <h2 className={styles.name}>{name}</h2>

        <div className={styles['price-wrapper']}>
          <strong className={styles.price}>{price}</strong>{' '}
          {interval && <div className={styles.interval}> / {interval}</div>}
        </div>

        <p className={styles.desc}>{desc}</p>

        <div className={styles['header-body']}>
          <div className={styles['cta-wrapper']}>{ctaContent}</div>
        </div>
      </div>

      {bottomContent}
    </div>
  );
};

export default Plan;
