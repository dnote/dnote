import React from 'react';
import classnames from 'classnames';

import Button from '../../Common/Button';
import ServerIcon from '../../Icons/Server';
import GlobeIcon from '../../Icons/Globe';

import styles from './Sidebar.module.scss';

const perks = [
  {
    id: 'hosted',
    icon: <ServerIcon width="16" height="16" fill="#4d4d8b" />,
    value: 'Fully hosted and managed'
  },
  {
    id: 'support',
    icon: <GlobeIcon width="16" height="16" fill="#4d4d8b" />,
    value: 'Support the Dnote community and development'
  }
];

function Sidebar({ isReady, transacting }) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.header}>
        <div className={styles['plan-name']}>Pro</div>

        <ul className={classnames('list-unstyled', styles.perks)}>
          {perks.map(perk => {
            return (
              <li key={perk.id} className={styles['perk-item']}>
                <div className={styles['perk-icon']}>{perk.icon}</div>
                <div className={styles['perk-value']}>{perk.value}</div>
              </li>
            );
          })}
        </ul>

        <div className={styles['price-wrapper']}>
          <strong className={styles.price}>$3.00</strong>
          <div className={styles.interval}>/ month</div>
        </div>

        <Button
          id="T-purchase-button"
          type="submit"
          className={classnames(
            'button button-large button-third button-stretch',
            styles['purchase-button']
          )}
          disabled={transacting}
          isBusy={transacting || !isReady}
        >
          Purchase Dnote Pro
        </Button>
      </div>

      <p className={styles.assurance}>You can cancel your plan any time.</p>
    </div>
  );
}

export default Sidebar;
