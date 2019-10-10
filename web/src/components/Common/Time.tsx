import React from 'react';

import {
  secondsToHTMLTimeDuration,
  getMonthName,
  getUTCOffset
} from '../../helpers/time';
import Tooltip from './Tooltip';
import { Alignment, Direction } from '../Common/Popover/types';
import styles from './Time.scss';

interface ContentProps {
  text: string;
  mobileText?: string;
}

const Content: React.FunctionComponent<ContentProps> = ({
  text,
  mobileText
}) => {
  if (mobileText === undefined) {
    return <span>{text}</span>;
  }

  return (
    <span>
      <span className={styles.text}>{text}</span>
      <span className={styles['mobile-text']}>{mobileText}</span>
    </span>
  );
};

interface Props {
  id: string;
  text: string;
  ms: number;
  mobileText?: string;
  isDuration?: boolean;
  wrapperClassName?: string;
  tooltipAlignment?: Alignment;
  tooltipDirection?: Direction;
}

function getDatetimeAttr(ms: number, isDuration: boolean = false): string {
  if (isDuration) {
    const s = ms / 1000;
    return secondsToHTMLTimeDuration(s);
  }

  const d = new Date(ms);
  return d.toISOString();
}

function formatOverlayTimeStr(ms: number): string {
  const date = new Date(ms);

  const y = date.getFullYear();
  const m = getMonthName(date, true);
  const d = date.getDate();
  const h = date.getHours();
  const min = date.getMinutes();
  const offset = getUTCOffset();

  let period;
  let hour;
  if (h >= 12) {
    period = 'PM';

    if (h === 12) {
      hour = h;
    } else {
      hour = h - 12;
    }
  } else {
    period = 'AM';
    hour = h;
  }

  return ` ${m} ${d}, ${y}, ${hour}:${min} ${period} GMT${offset}`;
}

const Time: React.FunctionComponent<Props> = ({
  id,
  text,
  mobileText,
  ms,
  isDuration,
  wrapperClassName,
  tooltipAlignment = 'center',
  tooltipDirection = 'bottom'
}) => {
  const dateTime = getDatetimeAttr(ms, isDuration);
  const overlay = <span>{formatOverlayTimeStr(ms)}</span>;

  return (
    <Tooltip
      id={id}
      alignment={tooltipAlignment}
      direction={tooltipDirection}
      overlay={overlay}
      wrapperClassName={wrapperClassName}
    >
      <time dateTime={dateTime} className={styles.time}>
        <Content text={text} mobileText={mobileText} />
      </time>
    </Tooltip>
  );
};

export default Time;
