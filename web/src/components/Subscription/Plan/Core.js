import React from 'react';

import BoxIcon from '../../Icons/Box';
import Plan from './internal';

const selfHostedPerks = [
  {
    id: 'own-machine',
    icon: <BoxIcon width="16" height="16" fill="#6e6e6e" />,
    value: 'Host on your own machine'
  }
];

function Core({ wrapperClassName, ctaContent, bottomContent }) {
  return (
    <Plan
      name="Core"
      price="Free"
      perks={selfHostedPerks}
      wrapperClassName={wrapperClassName}
      ctaContent={ctaContent}
      bottomContent={bottomContent}
    />
  );
}

export default Core;
