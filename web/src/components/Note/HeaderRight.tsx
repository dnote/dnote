import React from 'react';

import GlobeIcon from '../Icons/Globe';
import Tooltip from '../Common/Tooltip';

interface Props {
  isOwner: boolean;
  isPublic: boolean;
}

const HeaderRight: React.FunctionComponent<Props> = ({ isOwner, isPublic }) => {
  if (!isOwner) {
    return null;
  }
  if (!isPublic) {
    return null;
  }

  const publicTooltip = 'Anyone on the Internet can see this note.';

  return (
    <Tooltip
      id="note-public-indicator"
      alignment="right"
      direction="bottom"
      overlay={publicTooltip}
    >
      <GlobeIcon
        fill="#8c8c8c"
        width={16}
        height={16}
        ariaLabel={publicTooltip}
      />
    </Tooltip>
  );
};

export default HeaderRight;
