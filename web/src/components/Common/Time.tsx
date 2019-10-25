import React from 'react';

import {
  msToHTMLTimeDuration,
  getMonthName,
  getUTCOffset
} from '../../helpers/time';
import formatTime from '../../helpers/time/format';
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
    return msToHTMLTimeDuration(ms);
  }

  const d = new Date(ms);
  return d.toISOString();
}

function formatOverlayTimeStr(ms: number): string {
  const date = new Date(ms);

  return formatTime(date, '%MMM %DD, %YYYY, %hh:%mm %A GMT%Z');
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
