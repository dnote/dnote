import React from 'react';
import classnames from 'classnames';

import Menu, { MenuOption } from '../../Common/Menu';
import { Alignment, Direction } from '../../Common/Menu/types';
import styles from './SelectMenu.scss';

interface Props {
  defaultCurrentOptionIdx: number;
  options: MenuOption[];
  optRefs: any[];
  triggerText: string;
  isOpen: boolean;
  setIsOpen: (boolean) => void;
  triggerId: string;
  headerText: string;
  menuId: string;
  alignment: Alignment;
  direction: Direction;

  disabled?: boolean;
  wrapperClassName?: string;
}

const StatusFilter: React.FunctionComponent<Props> = ({
  defaultCurrentOptionIdx,
  options,
  optRefs,
  triggerText,
  disabled,
  isOpen,
  setIsOpen,
  headerText,
  triggerId,
  menuId,
  alignment,
  direction,
  wrapperClassName
}) => {
  return (
    <Menu
      defaultCurrentOptionIdx={defaultCurrentOptionIdx}
      options={options}
      disabled={disabled}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      optRefs={optRefs}
      triggerId={triggerId}
      menuId={menuId}
      triggerContent={
        <div className={styles['trigger-content']}>
          {triggerText}
          <span className="dropdown-caret" />
        </div>
      }
      headerContent={<div className={styles.header}>{headerText}</div>}
      triggerClassName={classnames(styles.trigger, {
        [styles['trigger-active']]: isOpen
      })}
      contentClassName={styles.content}
      wrapperClassName={wrapperClassName}
      alignment={alignment}
      direction={direction}
    />
  );
};

export default StatusFilter;
