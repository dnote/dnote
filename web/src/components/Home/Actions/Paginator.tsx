import React from 'react';
import classnames from 'classnames';

import { Link } from 'react-router-dom';
import { useFilters, useSelector } from '../../../store';
import { getHomePath } from '../../../libs/paths';
import CaretIcon from '../../Icons/Caret';
import styles from './Paginator.scss';

// PER_PAGE is the number of results per page in the response from the backend implementation's API.
// Currently it is fixed.
const PER_PAGE = 30;

type Direction = 'next' | 'prev';

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

interface PageLinkProps {
  direction: Direction;
  disabled: boolean;
  className?: string;
}

const PageLink: React.SFC<PageLinkProps> = ({
  direction,
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
      to={getHomePath({
        ...filters.queries,
        page
      })}
      aria-label={label}
      className={classnames(styles.link, className)}
    >
      {renderCaret(direction, 'black')}
    </Link>
  );
};

interface PaginatorProps {}

const Paginator: React.SFC<PaginatorProps> = () => {
  const filters = useFilters();
  const { notes } = useSelector(state => {
    return {
      notes: state.notes
    };
  });

  const hasNext = filters.page * PER_PAGE < notes.total;
  const hasPrev = filters.page > 1;
  const maxPage = Math.ceil(notes.total / PER_PAGE);

  let currentPage;
  if (maxPage > 0) {
    currentPage = filters.page;
  } else {
    currentPage = 0;
  }

  return (
    <nav className={styles.wrapper}>
      <span className={styles.info}>
        <span className={styles.label}>{currentPage}</span> of{' '}
        <span className={styles.label}>{maxPage}</span>
      </span>

      <PageLink
        direction="prev"
        disabled={!hasPrev}
        className={styles['link-prev']}
      />
      <PageLink direction="next" disabled={!hasNext} />
    </nav>
  );
};

export default Paginator;
