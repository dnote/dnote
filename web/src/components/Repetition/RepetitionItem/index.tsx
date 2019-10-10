import React, { useState } from 'react';
import classnames from 'classnames';

import { RepetitionRuleData } from 'jslib/operations/types';
import {
  secondsToDuration,
  secondsToHTMLTimeDuration,
  timeAgo
} from 'web/helpers/time';
import formatTime from 'web/helpers/time/format';
import Actions from './Actions';
import BookMeta from './BookMeta';
import Time from '../../Common/Time';
import styles from './RepetitionItem.scss';

interface Props {
  item: RepetitionRuleData;
  setRuleUUIDToDelete: React.Dispatch<any>;
}

function formatLastActive(ms: number): string {
  return timeAgo(ms);
}

const RepetitionItem: React.FunctionComponent<Props> = ({
  item,
  setRuleUUIDToDelete
}) => {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <li
      className={styles.wrapper}
      onMouseEnter={() => {
        setIsHovered(true);
      }}
      onMouseLeave={() => {
        setIsHovered(false);
      }}
    >
      <div className={styles.content}>
        <div className={styles.left}>
          <h2 className={styles.title}>{item.title}</h2>

          <div className={styles.meta}>
            <div>
              <span className={styles.frequency}>
                Every{' '}
                <time dateTime={secondsToHTMLTimeDuration(item.frequency)}>
                  {secondsToDuration(item.frequency)}
                </time>
              </span>
              <span className={styles.sep}>&middot;</span>
              <span className={styles.delivery}>email</span>
            </div>

            <BookMeta
              bookDomain={item.bookDomain}
              bookCount={item.books.length}
            />
          </div>
        </div>

        <div className={styles.right}>
          <ul className={classnames('list-unstyled', styles['detail-list'])}>
            <li>
              Last active:{' '}
              {item.lastActive === 0 ? (
                <span>Never</span>
              ) : (
                <Time
                  id={`${item.uuid}-lastactive-ts`}
                  text={formatLastActive(item.lastActive)}
                  ms={item.lastActive}
                  tooltipAlignment="center"
                  tooltipDirection="bottom"
                />
              )}
            </li>
            <li>
              Created:{' '}
              <Time
                id={`${item.uuid}-created-ts`}
                text={formatTime(new Date(item.createdAt), '%YYYY %MMM %Do')}
                ms={new Date(item.createdAt).getTime()}
                tooltipAlignment="center"
                tooltipDirection="bottom"
              />
            </li>
          </ul>
        </div>
      </div>

      <Actions
        isActive={isHovered}
        onDelete={() => {
          setRuleUUIDToDelete(item.uuid);
        }}
      />
    </li>
  );
};

export default RepetitionItem;
