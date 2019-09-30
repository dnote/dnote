import React from 'react';

import styles from './Preview.scss';
import { parseMarkdown } from '../../../helpers/markdown';

interface Props {
  content: string;
}

const Preview: React.SFC<Props> = ({ content }) => {
  return (
    <div
      className={styles.wrapper}
      // eslint-disable-next-line react/no-danger
      dangerouslySetInnerHTML={{ __html: parseMarkdown(content) }}
    />
  );
};

export default Preview;
