import React from 'react';
import classnames from 'classnames';

import { countries } from '../../../libs/countries';
import CaretIcon from '../../Icons/Caret';

import styles from './CountrySelect.module.scss';

function CountrySelect({ id, className, onChange, value }) {
  return (
    <div className={styles.wrapper}>
      <select
        id={id}
        className={classnames(className, styles.select)}
        value={value}
        onChange={onChange}
      >
        <option value="" />

        {countries.map(country => {
          return (
            <option key={country.code} value={country.code}>
              {country.name}
            </option>
          );
        })}
      </select>

      <CaretIcon width="12" height="12" className={styles.caret} />
    </div>
  );
}

export default CountrySelect;
