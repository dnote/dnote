import React from 'react';

const Icon = ({ fill, width, height, className }) => {
  const h = `${height}px`;
  const w = `${width}px`;

  return (
    <svg
      viewBox="0 0 32 32"
      height={h}
      width={w}
      fill={fill}
      className={className}
    >
      <g transform="translate(240 0)">
        <path d="M-211,4v26h-24c-1.104,0-2-0.895-2-2s0.896-2,2-2h22V0h-22c-2.209,0-4,1.791-4,4v24c0,2.209,1.791,4,4,4h26V4H-211z    M-235,8V2h20v22h-20V8z M-219,6h-12V4h12V6z M-223,10h-8V8h8V10z M-227,14h-4v-2h4V14z" />
      </g>
    </svg>
  );
};

Icon.defaultProps = {
  fill: '#000',
  width: 32,
  height: 32
};

export default Icon;
