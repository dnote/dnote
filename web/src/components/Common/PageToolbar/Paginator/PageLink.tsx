import React from 'react';
import classnames from 'classnames';
import { Location } from 'history';
import { Link } from 'react-router-dom';

import CaretIcon from '../../../Icons/Caret';
import { useFilters } from '../../../../store';
import styles from './Paginator.scss';

type Direction = 'next' | 'prev';

interface Props {
  direction: Direction;
  disabled: boolean;
  getPath: (page: number) => Location;
  className?: string;
}

const renderCaret = (direction: Direction, fill: string) => {
  return (
    <CaretIcon
      fill={fill}
      width={12}
      height={12}
      className={styles[`caret-${direction}`]}
    />
  );
};

const PageLink: React.FunctionComponent<Props> = ({
  direction,
  getPath,
  disabled,
  className
}) => {
  const filters = useFilters();

  if (disabled) {
    return (
      <span className={classnames(styles.link, styles.disabled, className)}>
        {renderCaret(direction, 'gray')}
      </span>
    );
  }

  let page;
  if (direction === 'next') {
    page = filters.page + 1;
  } else {
    page = filters.page - 1;
  }

  let label;
  if (direction === 'next') {
    label = 'Next page';
  } else {
    label = 'Previous page';
  }

  return (
    <Link
      to={getPath(page)}
      aria-label={label}
      className={classnames(styles.link, className)}
    >
      {renderCaret(direction, 'black')}
    </Link>
  );
};

export default PageLink;
