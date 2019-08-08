import React from 'react';
import classnames from 'classnames';
import { withRouter, RouteComponentProps } from 'react-router-dom';

import styles from './TabBar.scss';
import Item from './Item';
import NoteIcon from '../Icons/Note';
import BookIcon from '../Icons/Book';
import DashboardIcon from '../Icons/Dashboard';
import DotsIcon from '../Icons/Dots';
import HomeIcon from '../Icons/Home';

interface Props extends RouteComponentProps<any> {}

const TabBar: React.SFC<Props> = ({ location }) => {
  return (
    <nav className={styles.wrapper}>
      <ul className={classnames(styles.list, 'list-unstyled')}>
        <Item
          to="/"
          label="Home"
          renderIcon={fill => <HomeIcon width={20} height={20} fill={fill} />}
          active={location.pathname === '/'}
        />
        <Item
          to="/books"
          label="Books"
          renderIcon={fill => <BookIcon width={20} height={20} fill={fill} />}
          active={location.pathname === '/books'}
        />
        <Item
          to="/random"
          label="Random"
          renderIcon={fill => (
            <DashboardIcon width={20} height={20} fill={fill} />
          )}
          active={location.pathname === '/digests'}
        />
        <Item
          to="/new"
          label="New"
          renderIcon={fill => <NoteIcon width={20} height={20} fill={fill} />}
          active={location.pathname === '/new'}
        />
        <Item
          to="/settings"
          label="More"
          renderIcon={fill => <DotsIcon width={20} height={20} fill={fill} />}
          active={location.pathname === '/settings'}
        />
      </ul>
    </nav>
  );
};

export default withRouter(TabBar);
