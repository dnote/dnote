import React, { useState, useRef } from 'react';
import { Link } from 'react-router-dom';

import { getDigestPath } from 'web/libs/paths';
import SelectMenu from '../../Common/PageToolbar/SelectMenu';
import selectMenuStyles from '../../Common/PageToolbar/SelectMenu.scss';
import { Sort } from '../types';

interface Props {
  digestUUID: string;
  sort: Sort;
  disabled?: boolean;
}

const SortMenu: React.FunctionComponent<Props> = ({
  digestUUID,
  sort,
  disabled
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const optRefs = [useRef(null), useRef(null)];

  const options = [
    {
      name: 'newest',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestPath(digestUUID)}
          onClick={() => {
            setIsOpen(false);
          }}
          ref={optRefs[0]}
          tabIndex={-1}
        >
          Newest
        </Link>
      )
    },
    {
      name: 'oldest',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestPath(digestUUID, { sort: Sort.Oldest })}
          onClick={() => {
            setIsOpen(false);
          }}
          ref={optRefs[1]}
          tabIndex={-1}
        >
          Oldest
        </Link>
      )
    }
  ];

  let defaultCurrentOptionIdx: number;
  let sortText: string;
  if (sort === Sort.Oldest) {
    defaultCurrentOptionIdx = 1;
    sortText = 'Oldest';
  } else {
    defaultCurrentOptionIdx = 0;
    sortText = 'Newest';
  }

  return (
    <SelectMenu
      defaultCurrentOptionIdx={defaultCurrentOptionIdx}
      options={options}
      disabled={disabled}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      optRefs={optRefs}
      triggerId="sort-menu-trigger"
      menuId="sort-menu"
      headerText="Sort by"
      triggerText={`Sort: ${sortText}`}
      alignment="right"
      direction="bottom"
    />
  );
};

export default SortMenu;
