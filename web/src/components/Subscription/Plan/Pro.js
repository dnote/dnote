import React from 'react';

import Plan from './internal';
import ServerIcon from '../../Icons/Server';
import GlobeIcon from '../../Icons/Globe';

const proPerks = [
  {
    id: 'hosted',
    icon: <ServerIcon width="16" height="16" fill="#4d4d8b" />,
    value: 'Fully hosted and managed'
  },
  {
    id: 'support',
    icon: <GlobeIcon width="16" height="16" fill="#4d4d8b" />,
    value: 'Support the Dnote community and development'
  }
];

function ProPlan({ wrapperClassName, ctaContent, bottomContent }) {
  return (
    <Plan
      name="Pro"
      price="$3"
      interval="month"
      perks={proPerks}
      wrapperClassName={wrapperClassName}
      ctaContent={ctaContent}
      bottomContent={bottomContent}
    />
  );
}

export default ProPlan;
