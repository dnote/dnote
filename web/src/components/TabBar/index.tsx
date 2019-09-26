import React from 'react';
import classnames from 'classnames';
import { withRouter, RouteComponentProps, Link } from 'react-router-dom';

import styles from './TabBar.scss';
import Item from './Item';
import NoteIcon from '../Icons/Note';
import BookIcon from '../Icons/Book';
// import DashboardIcon from '../Icons/Dashboard';
import DotsIcon from '../Icons/Dots';
import HomeIcon from '../Icons/Home';

interface Props extends RouteComponentProps<any> {
  isMobileMenuOpen: boolean;
  setMobileMenuOpen: (boolean) => void;
}

function getFill(active: boolean): string {
  let ret;
  if (active) {
    ret = '#49abfd';
  } else {
    ret = '#cecece';
  }

  return ret;
}

const TabBar: React.SFC<Props> = ({
  location,
  isMobileMenuOpen,
  setMobileMenuOpen
}) => {
  const isHomeActive = !isMobileMenuOpen && location.pathname === '/';
  const isBookActive = !isMobileMenuOpen && location.pathname === '/books';
  // const isRandomActive = !isMobileMenuOpen && location.pathname === '/random';
  const isNewActive = !isMobileMenuOpen && location.pathname === '/new';

  return (
    <nav className={styles.wrapper}>
      <ul className={classnames(styles.list, 'list-unstyled')}>
        <Item>
          <Link
            to="/"
            className={classnames(styles.link, {
              [styles.active]: isHomeActive
            })}
          >
            <HomeIcon width={16} height={16} fill={getFill(isHomeActive)} />
            <span className={styles.label}>Home</span>
          </Link>
        </Item>

        <Item>
          <Link
            to="/books"
            className={classnames(styles.link, {
              [styles.active]: isBookActive
            })}
          >
            <BookIcon width={16} height={16} fill={getFill(isBookActive)} />
            <span className={styles.label}>Books</span>
          </Link>
        </Item>

        {/*
        <Item>
          <Link
            to="/random"
            className={classnames(styles.link, {
              [styles.active]: isRandomActive
            })}
          >
            <DashboardIcon
              width={16}
              height={16}
              fill={getFill(isRandomActive)}
            />
            <span className={styles.label}>Random</span>
          </Link>
        </Item>
        */}

        <Item>
          <Link
            to="/new"
            className={classnames(styles.link, {
              [styles.active]: isNewActive
            })}
          >
            <NoteIcon width={16} height={16} fill={getFill(isNewActive)} />
            <span className={styles.label}>New</span>
          </Link>
        </Item>

        <Item>
          <button
            type="button"
            className={classnames(styles.link, 'button-no-ui', {
              [styles.active]: isMobileMenuOpen
            })}
            onClick={() => {
              setMobileMenuOpen(!isMobileMenuOpen);
            }}
          >
            <DotsIcon width={16} height={16} fill={getFill(isMobileMenuOpen)} />
            <span className={styles.label}>More</span>
          </button>
        </Item>
      </ul>
    </nav>
  );
};

export default withRouter(TabBar);
