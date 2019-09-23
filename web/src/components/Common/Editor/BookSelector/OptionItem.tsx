import React from 'react';
import classnames from 'classnames';

import { Option } from 'jslib/helpers/select';
import CheckIcon from '../../../Icons/Check';
import styles from './OptionItem.scss';

interface Props {
  option: Option;
  isSelected: boolean;
  isFocused: boolean;
  onSelect: (Option) => void;
  isNew?: boolean;
  id?: string;
  className?: string;
}

function renderBody(isNew: boolean, label: string) {
  if (isNew) {
    return `Create book '${label}'`;
  }

  return label;
}

const OptionItem: React.SFC<Props> = ({
  option,
  isSelected,
  isFocused,
  onSelect,
  isNew,
  id
}) => {
  return (
    <button
      id={id}
      role="option"
      type="button"
      aria-selected={isSelected}
      onClick={() => {
        onSelect(option);
      }}
      className={classnames(
        'T-book-item-option',
        'button-no-ui',
        `book-item-${option.value}`,
        styles['combobox-option'],
        {
          [styles.active]: isSelected,
          [styles.focused]: isFocused
        }
      )}
    >
      {isSelected && (
        <CheckIcon
          fill="white"
          width={12}
          height={12}
          className={styles['check-icon']}
        />
      )}
      <span className={styles['option-label']}>
        {renderBody(isNew, option.label)}
      </span>
    </button>
  );
};

export default OptionItem;
