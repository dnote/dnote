import React from 'react';
import classnames from 'classnames';

import Button from '../../Common/Button';
import styles from './ShareModal.scss';

interface Props {
  kind: string;
  size?: string;
  isHot: boolean;
  onClick: () => void;
  className?: string;
}

const CopyButton: React.FunctionComponent<Props> = ({
  kind,
  size,
  isHot,
  onClick,
  className
}) => {
  return (
    <Button
      type="button"
      size={size}
      kind={kind}
      onClick={onClick}
      disabled={isHot}
      className={className}
    >
      {isHot ? 'Copied' : 'Copy link'}
    </Button>
  );
};

export default CopyButton;
