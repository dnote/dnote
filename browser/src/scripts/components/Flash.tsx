import React from 'react';

interface Props {
  message: string;
  when: boolean;
}

const Flash: React.FunctionComponent<Props> = ({ message, when }) => {
  if (when) {
    return <div className="alert error">Error: {message}</div>;
  }

  return null;
};

export default Flash;
